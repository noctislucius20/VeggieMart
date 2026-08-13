package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"order-service/config"
	"order-service/internal/adapter/repository"
	"order-service/internal/core/domain/entity"
	"order-service/internal/core/service/transaction"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderConsumerWorkerInterface interface {
	StartCreateOrderWorker(ctx context.Context)
	StartOrderPaymentSuccessWorker(ctx context.Context)
	StartElasticUpdateStatusOrderWorker(ctx context.Context)
	StartDbUpdateStatusOrderWorker(ctx context.Context)
}

type orderConsumerWorker struct {
	conn      *amqp.Connection
	esClient  *elasticsearch.Client
	repoOrder repository.OrderRepositoryInterface
	txManager transaction.TransactionManager
	logger    *log.Logger
	cfg       *config.Config
}

func NewOrderConsumerWorker(conn *amqp.Connection, esClient *elasticsearch.Client, repoOrder repository.OrderRepositoryInterface, txManager transaction.TransactionManager, cfg *config.Config, logger *log.Logger) OrderConsumerWorkerInterface {
	return &orderConsumerWorker{
		conn:      conn,
		esClient:  esClient,
		repoOrder: repoOrder,
		txManager: txManager,
		cfg:       cfg,
		logger:    logger,
	}
}

// StartDbUpdateStatusOrderWorker implements [OrderConsumerWorkerInterface].
func (o *orderConsumerWorker) StartDbUpdateStatusOrderWorker(ctx context.Context) {
	ch, err := o.conn.Channel()
	if err != nil {
		o.logger.Errorf("[OrderConsumer] StartDbUpdateStatusOrderWorker: %v", err)
		return
	}

	defer ch.Close()

	orderUpdateStatus := o.cfg.PublisherName.DbOrderUpdateStatus

	queue, err := ch.QueueDeclare(orderUpdateStatus, true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer] StartDbUpdateStatusOrderWorker: %v", err)
		return
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer] StartDbUpdateStatusOrderWorker: %v", err)
		return
	}

	o.logger.Infof("[OrderConsumer] StartDbUpdateStatusOrderWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				o.logger.Infof("[OrderConsumer] StartDbUpdateStatusOrderWorker: %v", "channel closed")
				return
			}

			var orderStatus entity.OrderEntity

			if err := json.Unmarshal(d.Body, &orderStatus); err != nil {
				o.logger.Errorf("[OrderConsumer] StartDbUpdateStatusOrderWorker: %v", err)
				continue
			}

			if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
				if err := o.repoOrder.UpdateOrderStatus(txCtx, orderStatus); err != nil {
					return err
				}

				return nil
			}); err != nil {
				o.logger.Errorf("[OrderConsumer] StartDbUpdateStatusOrderWorker: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			o.logger.Infof("[OrderConsumer] StartUpdateStatusOrderWorker: order %d successfully updated to elasticsearch", orderStatus.ID)
		}
	}
}

// StartElasticUpdateStatusOrderWorker implements [OrderConsumerWorkerInterface].
func (o *orderConsumerWorker) StartElasticUpdateStatusOrderWorker(ctx context.Context) {
	ch, err := o.conn.Channel()
	if err != nil {
		o.logger.Errorf("[OrderConsumer] StartElasticUpdateStatusOrderWorker: %v", err)
		return
	}

	defer ch.Close()

	orderUpdateStatus := o.cfg.PublisherName.ElasticOrderUpdateStatus

	queue, err := ch.QueueDeclare(orderUpdateStatus, true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer] StartElasticUpdateStatusOrderWorker: %v", err)
		return
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer] StartElasticUpdateStatusOrderWorker: %v", err)
		return
	}

	o.logger.Infof("[OrderConsumer] StartElasticUpdateStatusOrderWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				o.logger.Infof("[OrderConsumer] StartElasticUpdateStatusOrderWorker: %v", "channel closed")
				return
			}

			var orderStatus struct {
				OrderID int64  `json:"id"`
				Status  string `json:"status"`
				Remarks string `json:"remarks"`
			}

			if err := json.Unmarshal(d.Body, &orderStatus); err != nil {
				o.logger.Errorf("[OrderConsumer] StartElasticUpdateStatusOrderWorker: %v", err)
				continue
			}

			reqBody := map[string]any{
				"doc": map[string]any{
					"status":  orderStatus.Status,
					"remarks": orderStatus.Remarks,
				},
			}

			orderStatusJson, err := json.Marshal(reqBody)
			if err != nil {
				o.logger.Errorf("[OrderConsumer] StartElasticUpdateStatusOrderWorker: %v", err)
				continue
			}

			if _, err := o.esClient.Update(
				"orders",
				fmt.Sprintf("%d", orderStatus.OrderID),
				bytes.NewReader(orderStatusJson),
				o.esClient.Update.WithContext(ctx),
			); err != nil {
				o.logger.Errorf("[OrderConsumer] StartElasticUpdateStatusOrderWorker: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			o.logger.Infof("[OrderConsumer] StartElasticUpdateStatusOrderWorker: order %d successfully updated to elasticsearch", orderStatus.OrderID)
		}
	}
}

// StartOrderPaymentSuccessWorker implements [OrderConsumerWorkerInterface].
func (o *orderConsumerWorker) StartOrderPaymentSuccessWorker(ctx context.Context) {
	ch, err := o.conn.Channel()
	if err != nil {
		o.logger.Errorf("[OrderConsumer] StartOrderPaymentSuccessWorker: %v", err)
		return
	}

	defer ch.Close()

	orderPaymentSuccess := o.cfg.PublisherName.OrderPaymentSuccess

	queue, err := ch.QueueDeclare(orderPaymentSuccess, true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer] StartOrderPaymentSuccessWorker: %v", err)
		return
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer] StartOrderPaymentSuccessWorker: %v", err)
		return
	}

	o.logger.Infof("[OrderConsumer] StartOrderPaymentSuccessWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				o.logger.Infof("[OrderConsumer] StartOrderPaymentSuccessWorker: %v", "channel closed")
				return
			}

			var payment struct {
				OrderID       int64  `json:"order_id"`
				PaymentMethod string `json:"payment_method"`
			}

			if err := json.Unmarshal(d.Body, &payment); err != nil {
				o.logger.Errorf("[OrderConsumer] StartOrderPaymentSuccessWorker: %v", err)
				continue
			}

			reqBody := map[string]any{
				"doc": map[string]string{
					"payment_method": payment.PaymentMethod,
				},
			}

			paymentJson, err := json.Marshal(reqBody)
			if err != nil {
				o.logger.Errorf("[OrderConsumer] StartOrderPaymentSuccessWorker: %v", err)
				continue
			}

			if _, err := o.esClient.Update(
				"orders",
				fmt.Sprintf("%d", payment.OrderID),
				bytes.NewReader(paymentJson),
				o.esClient.Update.WithContext(ctx),
			); err != nil {
				o.logger.Errorf("[OrderConsumer] StartOrderPaymentSuccessWorker: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// body, _ := io.ReadAll(res.Body)
			// defer res.Body.Close()

			o.logger.Infof("[OrderConsumer] StartOrderPaymentSuccessWorker: order %d successfully updated to elasticsearch", payment.OrderID)
		}
	}
}

// StartCreateOrderWorker implements [OrderConsumerWorkerInterface].
func (o *orderConsumerWorker) StartCreateOrderWorker(ctx context.Context) {
	ch, err := o.conn.Channel()
	if err != nil {
		o.logger.Errorf("[OrderConsumer] StartCreateOrderWorker: %v", err)
		return
	}

	defer ch.Close()

	orderCreate := o.cfg.PublisherName.OrderCreate

	queue, err := ch.QueueDeclare(orderCreate, true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer] StartCreateOrderWorker: %v", err)
		return
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer] StartCreateOrderWorker: %v", err)
		return
	}

	o.logger.Infof("[OrderConsumer] StartCreateOrderWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				o.logger.Infof("[OrderConsumer] StartCreateOrderWorker: %v", "channel closed")
				return
			}

			var order entity.OrderEntity

			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				o.logger.Errorf("[OrderConsumer] StartCreateOrderWorker: %v", err)
				continue
			}

			orderJson, err := json.Marshal(order)
			if err != nil {
				o.logger.Errorf("[OrderConsumer] StartCreateOrderWorker: %v", err)
				continue
			}

			if _, err := o.esClient.Index(
				"orders",
				bytes.NewReader(orderJson),
				o.esClient.Index.WithDocumentID(fmt.Sprintf("%d", order.ID)),
				o.esClient.Index.WithContext(ctx),
				o.esClient.Index.WithRefresh("true"),
			); err != nil {
				o.logger.Errorf("[OrderConsumer] StartCreateOrderWorker: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			o.logger.Infof("[OrderConsumer] StartCreateOrderWorker: order %d successfully indexed to elasticsearch", order.ID)
		}
	}
}
