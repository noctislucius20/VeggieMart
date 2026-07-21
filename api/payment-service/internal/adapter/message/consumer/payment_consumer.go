package consumer

import (
	"context"
	"encoding/json"
	"payment-service/config"
	"payment-service/internal/adapter/repository"
	"payment-service/internal/adapter/repository/cache"
	"payment-service/internal/core/domain/entity"
	"payment-service/internal/core/service/transaction"
	"strings"

	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type PaymentConsumerWorkerInterface interface {
	StartCreatePaymentWorker(ctx context.Context)
	StartUpdatePaymentWorker(ctx context.Context)
}

type paymentConsumerWorker struct {
	conn         *amqp.Connection
	db           *gorm.DB
	logger       *log.Logger
	cfg          *config.Config
	cachePayment cache.PaymentCacheInterface
	paymentRepo  repository.PaymentRepositoryInterface
	txManager    transaction.TransactionManager
}

func NewPaymentConsumerWorker(conn *amqp.Connection, db *gorm.DB, cfg *config.Config, cachePayment cache.PaymentCacheInterface, paymentRepo repository.PaymentRepositoryInterface, txManager transaction.TransactionManager, logger *log.Logger) PaymentConsumerWorkerInterface {
	return &paymentConsumerWorker{
		conn:         conn,
		db:           db,
		cfg:          cfg,
		cachePayment: cachePayment,
		paymentRepo:  paymentRepo,
		txManager:    txManager,
		logger:       logger,
	}
}

// StartUpdatePaymentWorker implements [PaymentConsumerWorkerInterface].
func (p *paymentConsumerWorker) StartUpdatePaymentWorker(ctx context.Context) {
	ch, err := p.conn.Channel()
	if err != nil {
		p.logger.Errorf("[PaymentConsumer-1] StartUpdatePaymentWorker: %v", err)
		return
	}

	defer ch.Close()

	paymentUpdate := p.cfg.PublisherName.PaymentUpdate

	queue, err := ch.QueueDeclare(paymentUpdate, true, false, false, false, nil)
	if err != nil {
		p.logger.Errorf("[PaymentConsumer-2] StartUpdatePaymentWorker: %v", err)
		return
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		p.logger.Errorf("[PaymentConsumer-3] StartUpdatePaymentWorker: %v", err)
		return
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		p.logger.Errorf("[PaymentConsumer-4] StartUpdatePaymentWorker: %v", err)
		return
	}

	p.logger.Infof("[PaymentConsumer-5] StartUpdatePaymentWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				p.logger.Infof("[PaymentConsumer-6] StartUpdatePaymentWorker: %v", "channel closed")
				d.Nack(false, true)
				continue
			}

			var payment struct {
				OrderID          int64  `json:"order_id"`
				PaymentGatewayID string `json:"payment_gateway_id"`
				PaymentURL       string `json:"payment_url"`
			}

			err := json.Unmarshal(d.Body, &payment)
			if err != nil {
				p.logger.Errorf("[PaymentConsumer-7] StartUpdatePaymentWorker: %v", err)
				d.Nack(false, false)
				continue
			}

			if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
				paymentId, err := p.paymentRepo.UpdatePaymentGatewayByOrderId(txCtx, payment.OrderID, &payment.PaymentGatewayID, &payment.PaymentURL)
				if err != nil {
					return err
				}

				if err := p.cachePayment.DeletePaymentCache(txCtx, int64(paymentId), payment.OrderID); err != nil {
					return err
				}

				return nil
			}); err != nil {
				d.Nack(false, true)
				p.logger.Errorf("[PaymentConsumer-8] StartUpdatePaymentWorker: %v", err)
				continue
			}

			d.Ack(false)

			p.logger.Infof("[PaymentConsumer-9] StartUpdatePaymentWorker: payment for order %d successfully updated", payment.OrderID)
		}
	}
}

// StartCreatePaymentWorker implements [PaymentConsumerWorkerInterface].
func (p *paymentConsumerWorker) StartCreatePaymentWorker(ctx context.Context) {
	ch, err := p.conn.Channel()
	if err != nil {
		p.logger.Errorf("[PaymentConsumer-1] StartCreatePaymentWorker: %v", err)
		return
	}

	defer ch.Close()

	orderPaymentCreate := p.cfg.PublisherName.OrderPaymentCreate

	queue, err := ch.QueueDeclare(orderPaymentCreate, true, false, false, false, nil)
	if err != nil {
		p.logger.Errorf("[PaymentConsumer-2] StartCreatePaymentWorker: %v", err)
		return
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		p.logger.Errorf("[PaymentConsumer-3] StartCreatePaymentWorker: %v", err)
		return
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		p.logger.Errorf("[PaymentConsumer-4] StartCreatePaymentWorker: %v", err)
		return
	}

	p.logger.Infof("[PaymentConsumer-5] StartCreatePaymentWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				p.logger.Infof("[PaymentConsumer-6] StartCreatePaymentWorker: %v", "channel closed")
				d.Nack(false, true)
				continue
			}

			var orderEntity entity.OrderEntity

			err := json.Unmarshal(d.Body, &orderEntity)
			if err != nil {
				p.logger.Errorf("[PaymentConsumer-7] StartCreatePaymentWorker: %v", err)
				d.Nack(false, false)
				continue
			}

			paymentEntity := &entity.PaymentEntity{
				Order: orderEntity,
			}

			switch strings.ToLower(orderEntity.PaymentMethod) {
			case "cod":
				paymentEntity.PaymentStatus = "SUCCESS"
			case "transfer":
				paymentEntity.PaymentStatus = "PENDING"
			default:
				p.logger.Errorf("[PaymentConsumer-8] StartCreatePaymentWorker: %v", "invalid payment method")
				d.Nack(false, false)
				continue
			}

			if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
				paymentId, paymentStatus, err := p.paymentRepo.CreatePayment(txCtx, paymentEntity)
				if err != nil {
					return err
				}

				if err := p.paymentRepo.CreatePaymentLog(txCtx, paymentId, paymentStatus); err != nil {
					return err
				}

				if err := p.cachePayment.DeletePaymentCache(txCtx, int64(paymentId), orderEntity.ID); err != nil {
					return err
				}

				paymentEntity.ID = paymentId

				return nil
			}); err != nil {
				d.Nack(false, true)
				p.logger.Errorf("[PaymentConsumer-9] StartCreatePaymentWorker: %v", err)
				continue
			}

			d.Ack(false)

			p.logger.Infof("[PaymentConsumer-10] StartCreatePaymentWorker: payment for order %s successfully created", orderEntity.OrderCode)
		}
	}
}
