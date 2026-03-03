package service

import (
	"context"
	"errors"
	"payment-service/config"
	httpclient "payment-service/internal/adapter/http_client"
	"payment-service/internal/adapter/repository"
	"payment-service/internal/core/domain/entity"
	"payment-service/utils"
	"strings"
	"sync"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type PaymentServiceInterface interface {
	ProcessPayment(ctx context.Context, payment entity.PaymentEntity, userData string) (*entity.PaymentEntity, error)
	UpdateStatusByOrderCode(ctx context.Context, orderCode string, status string) error
	GetAllPayments(ctx context.Context, query entity.QueryStringPayment, userData string) ([]entity.PaymentEntity, int64, int64, error)
	GetPaymentById(ctx context.Context, paymentId uint, jwtUserData entity.JwtUserData, userData string) (*entity.PaymentEntity, error)

	getAllOrders(ctx context.Context, payments []entity.PaymentEntity, userData string) error
	getByIdOrder(ctx context.Context, payment *entity.PaymentEntity, userData string) error
}

type paymentService struct {
	repo        repository.PaymentRepositoryInterface
	repoOutbox  repository.OutboxEventInterface
	httpService HttpServiceInterface
	midtrans    httpclient.MidtransClientInterface
	cfg         *config.Config
	db          *gorm.DB
	logger      *log.Logger
}

// GetPaymentById implements [PaymentServiceInterface].
func (p *paymentService) GetPaymentById(ctx context.Context, paymentId uint, jwtUserData entity.JwtUserData, userData string) (*entity.PaymentEntity, error) {
	payment := &entity.PaymentEntity{}

	if err := p.db.Transaction(func(tx *gorm.DB) error {
		if strings.ToLower(jwtUserData.RoleName) == "customer" {
			paymentEntity, err := p.repo.GetPaymentById(ctx, paymentId, uint(jwtUserData.UserID), tx)
			if err != nil {
				return err
			}

			if err := p.getByIdOrder(ctx, paymentEntity, userData); err != nil {
				return err
			}

			payment = paymentEntity

			return nil
		}
		paymentEntity, err := p.repo.GetPaymentById(ctx, paymentId, 0, tx)
		if err != nil {
			return err
		}

		if err := p.getByIdOrder(ctx, paymentEntity, userData); err != nil {
			return err
		}

		payment = paymentEntity

		return nil
	}); err != nil {
		p.logger.Errorf("[PaymentService-1] GetPaymentById: %v", err.Error())
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

	if err := p.db.Transaction(func(tx *gorm.DB) error {
		paymentEntities, count, pages, err := p.repo.GetAllPayments(ctx, query, tx)
		if err != nil {
			return err
		}

		if len(paymentEntities) == 0 {
			return nil
		}

		if err := p.getAllOrders(ctx, paymentEntities, userData); err != nil {
			return err
		}

		payments, countData, totalPages = paymentEntities, count, pages

		return nil
	}); err != nil {
		p.logger.Errorf("[PaymentService-1] GetAllPayments: %v", err.Error())
		return nil, 0, 0, err
	}

	return payments, countData, totalPages, nil
}

// UpdateStatusByOrderCode implements [PaymentServiceInterface].
func (p *paymentService) UpdateStatusByOrderCode(ctx context.Context, orderCode string, status string) error {
	if err := p.db.Transaction(func(tx *gorm.DB) error {
		orderId, err := p.httpService.HttpOrderIdByOrderCodePublicService(orderCode)
		if err != nil {
			return err
		}

		if err := p.repo.UpdateStatusByOrderCode(ctx, orderId, status, tx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[PaymentService-1] UpdateStatusByOrderCode: %v", err.Error())
		return err
	}

	return nil
}

// ProcessPayment implements [PaymentServiceInterface].
func (p *paymentService) ProcessPayment(ctx context.Context, payment entity.PaymentEntity, userData string) (*entity.PaymentEntity, error) {
	publishPaymentSuccess := p.cfg.PublisherName.PaymentSuccess

	if err := p.db.Transaction(func(tx *gorm.DB) error {
		switch strings.ToLower(payment.PaymentMethod) {
		case "cod":
			payment.PaymentStatus = "SUCCESS"

			paymentId, paymentStatus, err := p.repo.CreatePayment(ctx, &payment, tx)
			if err != nil {
				return err
			}

			if err := p.repo.CreatePaymentLog(ctx, paymentId, paymentStatus, tx); err != nil {
				return err
			}

			// paymentEntity, err := p.repo.GetPaymentById(ctx, paymentId, 0, tx)
			// if err != nil {
			// 	return err
			// }

			payloadPublish := map[string]any{
				"order_id":       payment.OrderID,
				"payment_method": payment.PaymentMethod,
			}

			if err := p.repoOutbox.CreateEvent(ctx, publishPaymentSuccess, payloadPublish, &paymentId, tx); err != nil {
				return err
			}

		case "midtrans":
			if err := p.getByIdOrder(ctx, &payment, userData); err != nil {
				return err
			}

			transactionId, err := p.midtrans.CreateTransaction(payment.Order.OrderCode, int64(payment.GrossAmount), payment.Customer.CustomerName, payment.Customer.CustomerEmail)
			if err != nil {
				return err
			}

			payment.PaymentStatus = "PENDING"
			payment.PaymentGatewayID = transactionId

			paymentId, paymentStatus, err := p.repo.CreatePayment(ctx, &payment, tx)
			if err != nil {
				return err
			}

			if err := p.repo.CreatePaymentLog(ctx, paymentId, paymentStatus, tx); err != nil {
				return err
			}

			// paymentEntity, err := p.repo.GetPaymentById(ctx, paymentId, 0, tx)
			// if err != nil {
			// 	return err
			// }

			payloadPublish := map[string]any{
				"order_id":       payment.OrderID,
				"payment_method": payment.PaymentMethod,
			}

			if err := p.repoOutbox.CreateEvent(ctx, publishPaymentSuccess, payloadPublish, &paymentId, tx); err != nil {
				return err
			}

		default:
			err := errors.New(utils.INVALID_PAYMENT_METHOD)
			return err
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[PaymentService-1] ProcessPayment: %v", err.Error())
		return nil, err
	}

	return &payment, nil
}

// getByIdOrder implements [PaymentServiceInterface].
func (p *paymentService) getByIdOrder(ctx context.Context, payment *entity.PaymentEntity, userData string) error {
	var (
		wg          sync.WaitGroup
		resultOrder *entity.OrderDetailResponseEntity
		errCh       = make(chan error, 1)
		err         error
	)

	wg.Go(func() {
		resultOrder, err = p.httpService.HttpOrderByIdService(int64(payment.OrderID), userData)
		if err != nil {
			errCh <- err
		}
	})

	wg.Wait()

	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	payment.Customer.CustomerName = resultOrder.Customer.CustomerName
	payment.Customer.CustomerEmail = resultOrder.Customer.CustomerEmail
	payment.Customer.CustomerAddress = resultOrder.Customer.CustomerAddress

	payment.Order.OrderCode = resultOrder.OrderCode
	payment.Order.OrderShippingType = resultOrder.ShippingType
	payment.Order.OrderAt = resultOrder.OrderDatetime
	payment.Order.OrderRemarks = resultOrder.Remarks

	return nil
}

// getAllOrders implements [PaymentServiceInterface].
func (p *paymentService) getAllOrders(ctx context.Context, payments []entity.PaymentEntity, userData string) error {
	orderIds := map[uint]struct{}{}
	for _, payment := range payments {
		orderIds[payment.OrderID] = struct{}{}
	}

	reqOrderIds := make([]int64, 0, len(orderIds))
	for id := range orderIds {
		reqOrderIds = append(reqOrderIds, int64(id))
	}

	var (
		wg           sync.WaitGroup
		resultOrders map[int64]entity.OrderDetailResponseEntity
		errCh        = make(chan error, 1)
		err          error
	)

	wg.Go(func() {
		resultOrders, err = p.httpService.HttpOrdersAllService(reqOrderIds, userData)
		if err != nil {
			errCh <- err
		}
	})

	wg.Wait()

	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	for pIdx, payment := range payments {
		if q, ok := resultOrders[int64(payment.OrderID)]; ok {
			payments[pIdx].Order.OrderCode = q.OrderCode
			payments[pIdx].Order.OrderShippingType = q.ShippingType
		}
	}

	return nil
}

func NewPaymentService(
	repo repository.PaymentRepositoryInterface,
	repoOutbox repository.OutboxEventInterface,
	cfg *config.Config,
	httpService HttpServiceInterface,
	midtrans httpclient.MidtransClientInterface,
	db *gorm.DB,
	logger *log.Logger,
) PaymentServiceInterface {
	return &paymentService{
		repo:        repo,
		repoOutbox:  repoOutbox,
		httpService: httpService,
		midtrans:    midtrans,
		cfg:         cfg,
		db:          db,
		logger:      logger,
	}
}
