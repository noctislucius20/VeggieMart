package service

import (
	"context"
	"errors"
	"payment-service/config"
	httpclient "payment-service/internal/adapter/http_client"
	"payment-service/internal/adapter/repository"
	"payment-service/internal/adapter/repository/cache"
	"payment-service/internal/core/domain/entity"
	"payment-service/internal/core/service/transaction"
	"payment-service/utils"
	"strings"

	"github.com/labstack/gommon/log"
)

type PaymentServiceInterface interface {
	ProcessPayment(ctx context.Context, payment entity.PaymentEntity, jwtUserData entity.JwtUserData) (*entity.PaymentEntity, error)
	UpdateStatusByOrderCode(ctx context.Context, orderCode string, status string) error
	GetAllPayments(ctx context.Context, query entity.QueryStringPayment, userData string) ([]entity.PaymentEntity, int64, int64, error)
	GetPaymentById(ctx context.Context, paymentId int64, jwtUserData entity.JwtUserData, userData string) (*entity.PaymentEntity, error)

	getOrderHttp(ctx context.Context, payment *entity.PaymentEntity, jwtUserData entity.JwtUserData, roleName string) error
}

type paymentService struct {
	repo         repository.PaymentRepositoryInterface
	repoOutbox   repository.OutboxEventInterface
	cachePayment cache.PaymentCacheInterface
	httpService  HttpServiceInterface
	midtrans     httpclient.MidtransClientInterface
	txManager    transaction.TransactionManager
	cfg          *config.Config
	logger       *log.Logger
}

func NewPaymentService(repo repository.PaymentRepositoryInterface, repoOutbox repository.OutboxEventInterface, cachePayment cache.PaymentCacheInterface, cfg *config.Config, httpService HttpServiceInterface, midtrans httpclient.MidtransClientInterface, txManager transaction.TransactionManager, logger *log.Logger) PaymentServiceInterface {
	return &paymentService{
		repo:         repo,
		repoOutbox:   repoOutbox,
		cachePayment: cachePayment,
		httpService:  httpService,
		midtrans:     midtrans,
		txManager:    txManager,
		cfg:          cfg,
		logger:       logger,
	}
}

// GetPaymentById implements [PaymentServiceInterface].
func (p *paymentService) GetPaymentById(ctx context.Context, paymentId int64, jwtUserData entity.JwtUserData, userData string) (*entity.PaymentEntity, error) {
	payment := &entity.PaymentEntity{}

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		roleEntity, err := p.cachePayment.GetRoleById(txCtx, jwtUserData.RoleID)
		if err != nil {
			return err
		}

		switch strings.ToLower(roleEntity.Name) {
		case "customer": // requested by customer
			paymentEntity, err := p.cachePayment.GetPaymentById(txCtx, paymentId, int64(jwtUserData.UserID))
			if err != nil {
				return err
			}

			payment = paymentEntity
		default: // requested by admin
			paymentEntity, err := p.cachePayment.GetPaymentById(txCtx, paymentId, 0)
			if err != nil {
				return err
			}

			payment = paymentEntity
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[PaymentService-1] GetPaymentById: %v", err)
		return nil, err
	}

	return payment, nil
}

// GetAllPayments implements [PaymentServiceInterface].
func (p *paymentService) GetAllPayments(ctx context.Context, query entity.QueryStringPayment, userData string) ([]entity.PaymentEntity, int64, int64, error) {
	var (
		payments   []entity.PaymentEntity
		countData  int64
		totalPages int64
	)

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		paymentEntities, count, pages, err := p.repo.GetAllPayments(txCtx, query)
		if err != nil {
			return err
		}

		if len(paymentEntities) == 0 {
			return nil
		}

		payments, countData, totalPages = paymentEntities, count, pages

		return nil
	}); err != nil {
		p.logger.Errorf("[PaymentService-1] GetAllPayments: %v", err)
		return nil, 0, 0, err
	}

	return payments, countData, totalPages, nil
}

// UpdateStatusByOrderCode implements [PaymentServiceInterface].
func (p *paymentService) UpdateStatusByOrderCode(ctx context.Context, orderCode string, status string) error {
	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		orderId, err := p.repo.GetOrderIdByOrderCode(txCtx, orderCode)
		if err != nil {
			return err
		}

		paymentId, err := p.repo.UpdateStatusByOrderCode(txCtx, orderId, status)
		if err != nil {
			return err
		}

		if err := p.cachePayment.DeletePaymentCache(txCtx, int64(paymentId)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[PaymentService-1] UpdateStatusByOrderCode: %v", err)
		return err
	}

	return nil
}

// ProcessPayment implements [PaymentServiceInterface].
func (p *paymentService) ProcessPayment(ctx context.Context, payment entity.PaymentEntity, jwtUserData entity.JwtUserData) (*entity.PaymentEntity, error) {
	publishPaymentSuccess := p.cfg.PublisherName.PaymentSuccess

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		roleEntity, err := p.cachePayment.GetRoleById(txCtx, jwtUserData.RoleID)
		if err != nil {
			return err
		}

		if err := p.getOrderHttp(txCtx, &payment, jwtUserData, roleEntity.Name); err != nil {
			return err
		}

		switch strings.ToLower(payment.PaymentMethod) {
		case "cod":
			payment.PaymentStatus = "SUCCESS"

			paymentId, paymentStatus, err := p.repo.CreatePayment(txCtx, &payment)
			if err != nil {
				return err
			}

			if err := p.cachePayment.DeletePaymentCache(txCtx, int64(paymentId)); err != nil {
				return err
			}

			if err := p.repo.CreatePaymentLog(txCtx, paymentId, paymentStatus); err != nil {
				return err
			}

			// paymentEntity, err := p.repo.GetPaymentById(txCtx, paymentId, 0)
			// if err != nil {
			// 	return err
			// }

			payloadPublish := map[string]any{
				"order_id":       payment.Order.ID,
				"payment_method": payment.PaymentMethod,
			}

			if err := p.repoOutbox.CreateEvent(txCtx, publishPaymentSuccess, payloadPublish, &paymentId); err != nil {
				return err
			}

			payment.ID = paymentId
		case "midtrans":
			transactionId, err := p.midtrans.CreateTransaction(payment.Order.OrderCode, int64(payment.GrossAmount), payment.Customer.CustomerName, payment.Customer.CustomerEmail)
			if err != nil {
				return err
			}

			payment.PaymentStatus = "PENDING"
			payment.PaymentGatewayID = transactionId

			paymentId, paymentStatus, err := p.repo.CreatePayment(txCtx, &payment)
			if err != nil {
				return err
			}

			if err := p.cachePayment.DeletePaymentCache(txCtx, int64(paymentId)); err != nil {
				return err
			}

			if err := p.repo.CreatePaymentLog(txCtx, paymentId, paymentStatus); err != nil {
				return err
			}

			// paymentEntity, err := p.repo.GetPaymentById(txCtx, paymentId, 0)
			// if err != nil {
			// 	return err
			// }

			payloadPublish := map[string]any{
				"order_id":       payment.Order.ID,
				"payment_method": payment.PaymentMethod,
			}

			if err := p.repoOutbox.CreateEvent(txCtx, publishPaymentSuccess, payloadPublish, &paymentId); err != nil {
				return err
			}

			payment.ID = paymentId
		default:
			return errors.New(utils.INVALID_PAYMENT_METHOD)
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[PaymentService-1] ProcessPayment: %v", err)
		return nil, err
	}

	return &payment, nil
}

// getOrderHttp implements [PaymentServiceInterface].
func (p *paymentService) getOrderHttp(ctx context.Context, payment *entity.PaymentEntity, jwtUserData entity.JwtUserData, roleName string) error {
	resultOrder, err := p.httpService.HttpOrderByIdService(int64(payment.Order.ID), jwtUserData, roleName)
	if err != nil {
		return err
	}

	payment.Customer.CustomerID = resultOrder.Customer.CustomerID
	payment.Customer.CustomerName = resultOrder.Customer.CustomerName
	payment.Customer.CustomerPhone = resultOrder.Customer.CustomerPhone
	payment.Customer.CustomerEmail = resultOrder.Customer.CustomerEmail
	payment.Customer.CustomerAddress = resultOrder.Customer.CustomerAddress

	payment.Order.ID = resultOrder.ID
	payment.Order.OrderCode = resultOrder.OrderCode
	payment.Order.OrderDatetime = resultOrder.OrderDatetime
	payment.Order.Remarks = resultOrder.Remarks
	payment.Order.ShippingFee = resultOrder.ShippingFee
	payment.Order.ShippingType = resultOrder.ShippingType
	payment.Order.Status = resultOrder.Status
	payment.Order.TotalAmount = resultOrder.TotalAmount

	return nil
}
