package service

import (
	"context"
	"encoding/json"
	"fmt"
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
	UpdateStatusByOrderCode(ctx context.Context, orderCode string, status string, customerEmail string) error
	GetAllPayments(ctx context.Context, query entity.QueryStringPayment, userData string) ([]entity.PaymentEntity, int64, int64, error)
	GetPaymentById(ctx context.Context, paymentId int64, jwtUserData entity.JwtUserData, userData string) (*entity.PaymentEntity, error)
	GetPaymentByOrderId(ctx context.Context, orderId int64, jwtUserData entity.JwtUserData) (*entity.PaymentEntity, error)

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
		p.logger.Errorf("[PaymentService] GetPaymentById: %v", err)
		return nil, err
	}

	return payment, nil
}

// GetPaymentByOrderId implements [PaymentServiceInterface].
func (p *paymentService) GetPaymentByOrderId(ctx context.Context, orderId int64, jwtUserData entity.JwtUserData) (*entity.PaymentEntity, error) {
	payment := &entity.PaymentEntity{}

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		roleEntity, err := p.cachePayment.GetRoleById(txCtx, jwtUserData.RoleID)
		if err != nil {
			return err
		}

		switch strings.ToLower(roleEntity.Name) {
		case "customer": // requested by customer
			paymentEntity, err := p.cachePayment.GetPaymentByOrderId(txCtx, orderId, int64(jwtUserData.UserID))
			if err != nil {
				return err
			}

			payment = paymentEntity
		default: // requested by admin
			paymentEntity, err := p.cachePayment.GetPaymentByOrderId(txCtx, orderId, 0)
			if err != nil {
				return err
			}

			payment = paymentEntity
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[PaymentService] GetPaymentByOrderId: %v", err)
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
		p.logger.Errorf("[PaymentService] GetAllPayments: %v", err)
		return nil, 0, 0, err
	}

	return payments, countData, totalPages, nil
}

// UpdateStatusByOrderCode implements [PaymentServiceInterface].
func (p *paymentService) UpdateStatusByOrderCode(ctx context.Context, orderCode string, status string, customerEmail string) error {
	var (
		publishEmailUpdateStatus        = p.cfg.PublisherName.EmailUpdateOrderStatus
		publishOrderElasticUpdateStatus = p.cfg.PublisherName.ElasticOrderUpdateStatus
		publishOrderDbUpdateStatus      = p.cfg.PublisherName.DbOrderUpdateStatus
		outboxEventEntities             []entity.OutboxEventEntity
	)
	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		paymentId, orderId, customerId, err := p.repo.GetPaymentOrderIdByOrderCode(txCtx, orderCode)
		if err != nil {
			return err
		}

		if err := p.repo.UpdateStatusByPaymentId(txCtx, paymentId, status); err != nil {
			return err
		}

		if err := p.cachePayment.DeletePaymentCache(txCtx, paymentId, orderId); err != nil {
			return err
		}

		if status == "SUCCESS" {
			// ADMIN_PUSH notification for transfer success
			adminPushMessage := fmt.Sprintf("Pembayaran pesanan %s senilai Rp telah sukses (Transfer).", orderCode)
			adminNotifPayload := map[string]any{
				"receiver_email":       "",
				"message":              adminPushMessage,
				"subject":              "Pembayaran Transfer Sukses",
				"notification_type":    "orders",
				"notification_type_id": orderId,
				"receiver_id":          0,
				"notification_method":  "ADMIN_PUSH",
			}

			jsonAdminNotif, _ := json.Marshal(adminNotifPayload)

			outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
				EventType:   utils.NOTIF_PUSH,
				Payload:     string(jsonAdminNotif),
				AggregateID: fmt.Sprintf("%d", paymentId),
			})

			if err := p.repoOutbox.CreateBatchEvents(txCtx, outboxEventEntities); err != nil {
				return err
			}
		}

		if status == "FAILED" {
			payloadMessage := fmt.Sprintf("Hello,\n\nYour order with ID %s has been updated with status: %s", orderCode, "CANCELLED")

			// consumed by notification service
			publishEmailPayload := map[string]any{
				"receiver_email":       customerEmail,
				"message":              payloadMessage,
				"subject":              "Update Status Order",
				"notification_type":    "orders",
				"notification_type_id": orderId,
				"receiver_id":          customerId,
				"notification_method":  "EMAIL",
			}

			jsonEmailUpdateStatus, _ := json.Marshal(publishEmailPayload)

			outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
				EventType:   publishEmailUpdateStatus,
				Payload:     string(jsonEmailUpdateStatus),
				AggregateID: fmt.Sprintf("%d", paymentId),
			})

			// consumed by notification service
			publishPushNotifPayload := map[string]any{
				"receiver_email":       "",
				"message":              payloadMessage,
				"subject":              "Update Status Order",
				"notification_type":    "orders",
				"notification_type_id": orderId,
				"receiver_id":          customerId,
				"notification_method":  "PUSH",
			}

			jsonPushNotif, _ := json.Marshal(publishPushNotifPayload)

			outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
				EventType:   utils.NOTIF_PUSH,
				Payload:     string(jsonPushNotif),
				AggregateID: fmt.Sprintf("%d", paymentId),
			})

			// consumed by order elastic and db
			publishOrderUpdateStatusPayload := map[string]any{
				"id":     orderId,
				"status": "CANCELLED",
			}

			jsonOrderUpdateStatus, _ := json.Marshal(publishOrderUpdateStatusPayload)

			outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
				EventType:   publishOrderElasticUpdateStatus,
				Payload:     string(jsonOrderUpdateStatus),
				AggregateID: fmt.Sprintf("%d", paymentId),
			})

			outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
				EventType:   publishOrderDbUpdateStatus,
				Payload:     string(jsonOrderUpdateStatus),
				AggregateID: fmt.Sprintf("%d", paymentId),
			})

			if err := p.repoOutbox.CreateBatchEvents(txCtx, outboxEventEntities); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[PaymentService] UpdateStatusByOrderCode: %v", err)
		return err
	}

	return nil
}

// ProcessPayment implements [PaymentServiceInterface].
func (p *paymentService) ProcessPayment(ctx context.Context, payment entity.PaymentEntity, jwtUserData entity.JwtUserData) (*entity.PaymentEntity, error) {
	var (
		publishPaymentSuccess = p.cfg.PublisherName.PaymentSuccess
		publishPaymentUpdate  = p.cfg.PublisherName.PaymentUpdate
		outboxEventEntities   []entity.OutboxEventEntity
	)

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
			// consumed by order service elastic
			payloadPublish := map[string]any{
				"order_id":       payment.Order.ID,
				"payment_method": payment.PaymentMethod,
			}

			if err := p.repoOutbox.CreateEvent(txCtx, publishPaymentSuccess, payloadPublish, &payment.Order.ID); err != nil {
				return err
			}

			// ADMIN_PUSH notification for COD success
			adminPushMessage := fmt.Sprintf("Pembayaran pesanan %s senilai Rp %d telah sukses (COD).", payment.Order.OrderCode, payment.Order.TotalAmount)
			adminNotifPayload := map[string]any{
				"receiver_email":       "",
				"message":              adminPushMessage,
				"subject":              "Pembayaran COD Sukses",
				"notification_type":    "orders",
				"notification_type_id": payment.Order.ID,
				"receiver_id":          0,
				"notification_method":  "ADMIN_PUSH",
			}

			jsonAdminNotif, _ := json.Marshal(adminNotifPayload)

			outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
				EventType:   utils.NOTIF_PUSH,
				Payload:     string(jsonAdminNotif),
				AggregateID: fmt.Sprintf("%d", payment.Order.ID),
			})

			if err := p.repoOutbox.CreateBatchEvents(txCtx, outboxEventEntities); err != nil {
				return err
			}
		case "transfer":
			transactionId, redirectUrl, err := p.midtrans.CreateTransaction(payment.Order.OrderCode, int64(payment.Order.TotalAmount), payment.Customer.CustomerName, payment.Customer.CustomerEmail)
			if err != nil {
				return err
			}

			payment.PaymentGatewayID = &transactionId
			payment.PaymentURL = &redirectUrl

			// consumed by order service elastic
			jsonPaymentSuccess, _ := json.Marshal(map[string]any{
				"order_id":       payment.Order.ID,
				"payment_method": payment.PaymentMethod,
			})

			outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
				EventType:   publishPaymentSuccess,
				Payload:     string(jsonPaymentSuccess),
				AggregateID: fmt.Sprintf("%d", payment.Order.ID),
			})

			// consumed by payment service db
			jsonUpdatePayment, _ := json.Marshal(map[string]any{
				"order_id":           payment.Order.ID,
				"payment_gateway_id": transactionId,
				"payment_url":        redirectUrl,
			})

			outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
				EventType:   publishPaymentUpdate,
				Payload:     string(jsonUpdatePayment),
				AggregateID: fmt.Sprintf("%d", payment.Order.ID),
			})

			if err := p.repoOutbox.CreateBatchEvents(txCtx, outboxEventEntities); err != nil {
				return err
			}

		default:
			return utils.ErrInvalidPaymentMethod
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[PaymentService] ProcessPayment: %v", err)
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
