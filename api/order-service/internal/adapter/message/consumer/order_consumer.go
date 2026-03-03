package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"order-service/config"
	"order-service/internal/core/domain/entity"

	"github.com/labstack/gommon/log"
	"github.com/rabbitmq/amqp091-go"
)

type OrderConsumerWorkerInterface interface {
	StartCreateOrderWorker(ctx context.Context)
	StartOrderPaymentSuccessWorker(ctx context.Context)
	StartUpdateStatusOrderWorker(ctx context.Context)
}

type orderConsumerWorker struct {
	conn   *amqp091.Connection
	logger *log.Logger
}

// StartUpdateStatusOrderWorker implements [OrderConsumerWorkerInterface].
func (o *orderConsumerWorker) StartUpdateStatusOrderWorker(ctx context.Context) {
	ch, err := o.conn.Channel()
	if err != nil {
		o.logger.Errorf("[OrderConsumer-1] StartUpdateStatusOrderWorker: %v", err.Error())
		return
	}

	defer ch.Close()

	orderUpdateStatus := config.NewConfig().PublisherName.OrderUpdateStatus

	queue, err := ch.QueueDeclare(orderUpdateStatus, true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer-2] StartUpdateStatusOrderWorker: %v", err.Error())
		return
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer-3] StartUpdateStatusOrderWorker: %v", err.Error())
		return
	}

	esClient, err := config.NewConfig().NewElasticsearchClient()
	if err != nil {
		o.logger.Errorf("[OrderConsumer-4] StartUpdateStatusOrderWorker: %v", err.Error())
		return
	}

	o.logger.Infof("[OrderConsumer-5] StartUpdateStatusOrderWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				o.logger.Infof("[OrderConsumer-6] StartUpdateStatusOrderWorker: %v", "channel closed")
				continue
			}

			var orderStatus struct {
				OrderID int64  `json:"id"`
				Status  string `json:"status"`
				Remarks string `json:"remarks"`
			}

			if err := json.Unmarshal(d.Body, &orderStatus); err != nil {
				o.logger.Errorf("[OrderConsumer-7] StartUpdateStatusOrderWorker: %v", err.Error())
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
				o.logger.Errorf("[OrderConsumer-8] StartUpdateStatusOrderWorker: %v", err.Error())
				continue
			}

			if _, err := esClient.Update(
				"orders",
				fmt.Sprintf("%d", orderStatus.OrderID),
				bytes.NewReader(orderStatusJson),
				esClient.Update.WithContext(ctx),
			); err != nil {
				o.logger.Errorf("[OrderConsumer-9] StartUpdateStatusOrderWorker: %v", err.Error())
				continue
			}

			// body, _ := io.ReadAll(res.Body)
			// defer res.Body.Close()

			o.logger.Infof("[OrderConsumer-10] StartUpdateStatusOrderWorker: order %d successfully updated to elasticsearch", orderStatus.OrderID)
		}
	}
}

// StartOrderPaymentSuccessWorker implements [OrderConsumerWorkerInterface].
func (o *orderConsumerWorker) StartOrderPaymentSuccessWorker(ctx context.Context) {
	ch, err := o.conn.Channel()
	if err != nil {
		o.logger.Errorf("[OrderConsumer-1] StartOrderPaymentSuccessWorker: %v", err.Error())
		return
	}

	defer ch.Close()

	orderPaymentSuccess := config.NewConfig().PublisherName.OrderPaymentSuccess

	queue, err := ch.QueueDeclare(orderPaymentSuccess, true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer-2] StartOrderPaymentSuccessWorker: %v", err.Error())
		return
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer-3] StartOrderPaymentSuccessWorker: %v", err.Error())
		return
	}

	esClient, err := config.NewConfig().NewElasticsearchClient()
	if err != nil {
		o.logger.Errorf("[OrderConsumer-4] StartOrderPaymentSuccessWorker: %v", err.Error())
		return
	}

	o.logger.Infof("[OrderConsumer-5] StartOrderPaymentSuccessWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				o.logger.Infof("[OrderConsumer-6] StartOrderPaymentSuccessWorker: %v", "channel closed")
				continue
			}

			var payment struct {
				OrderID       int64  `json:"order_id"`
				PaymentMethod string `json:"payment_method"`
			}

			if err := json.Unmarshal(d.Body, &payment); err != nil {
				o.logger.Errorf("[OrderConsumer-7] StartOrderPaymentSuccessWorker: %v", err.Error())
				continue
			}

			reqBody := map[string]any{
				"doc": map[string]string{
					"payment_method": payment.PaymentMethod,
				},
			}

			paymentJson, err := json.Marshal(reqBody)
			if err != nil {
				o.logger.Errorf("[OrderConsumer-8] StartOrderPaymentSuccessWorker: %v", err.Error())
				continue
			}

			if _, err := esClient.Update(
				"orders",
				fmt.Sprintf("%d", payment.OrderID),
				bytes.NewReader(paymentJson),
				esClient.Update.WithContext(ctx),
			); err != nil {
				o.logger.Errorf("[OrderConsumer-9] StartOrderPaymentSuccessWorker: %v", err.Error())
				continue
			}

			// body, _ := io.ReadAll(res.Body)
			// defer res.Body.Close()

			o.logger.Infof("[OrderConsumer-10] StartOrderPaymentSuccessWorker: order %d successfully updated to elasticsearch", payment.OrderID)
		}
	}
}

// StartCreateOrderWorker implements [OrderConsumerWorkerInterface].
func (o *orderConsumerWorker) StartCreateOrderWorker(ctx context.Context) {
	ch, err := o.conn.Channel()
	if err != nil {
		o.logger.Errorf("[OrderConsumer-1] StartCreateOrderWorker: %v", err.Error())
		return
	}

	defer ch.Close()

	orderCreate := config.NewConfig().PublisherName.OrderCreate

	queue, err := ch.QueueDeclare(orderCreate, true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer-2] StartCreateOrderWorker: %v", err.Error())
		return
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		o.logger.Errorf("[OrderConsumer-3] StartCreateOrderWorker: %v", err.Error())
		return
	}

	esClient, err := config.NewConfig().NewElasticsearchClient()
	if err != nil {
		o.logger.Errorf("[OrderConsumer-4] StartCreateOrderWorker: %v", err.Error())
		return
	}

	o.logger.Infof("[OrderConsumer-5] StartCreateOrderWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				o.logger.Infof("[OrderConsumer-6] StartCreateOrderWorker: %v", "channel closed")
				continue
			}

			var order entity.OrderEntity

			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				o.logger.Errorf("[OrderConsumer-7] StartCreateOrderWorker: %v", err.Error())
				continue
			}

			orderJson, err := json.Marshal(order)
			if err != nil {
				o.logger.Errorf("[OrderConsumer-8] StartCreateOrderWorker: %v", err.Error())
				continue
			}

			if _, err := esClient.Index(
				"orders",
				bytes.NewReader(orderJson),
				esClient.Index.WithDocumentID(fmt.Sprintf("%d", order.ID)),
				esClient.Index.WithContext(ctx),
				esClient.Index.WithRefresh("true"),
			); err != nil {
				o.logger.Errorf("[OrderConsumer-9] StartCreateOrderWorker: %v", err.Error())
				continue
			}

			// body, _ := io.ReadAll(res.Body)
			// defer res.Body.Close()

			o.logger.Infof("[OrderConsumer-10] StartCreateOrderWorker: order %d successfully indexed to elasticsearch", order.ID)
		}
	}
}

func NewOrderConsumerWorker(conn *amqp091.Connection, logger *log.Logger) OrderConsumerWorkerInterface {
	return &orderConsumerWorker{
		conn:   conn,
		logger: logger,
	}
}
