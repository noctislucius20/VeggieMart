package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"product-service/config"
	"product-service/internal/core/domain/entity"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ProductConsumerWorkerInterface interface {
	StartCreateProductWorker(ctx context.Context)
	StartDeleteProductWorker(ctx context.Context)
	StartUpdateProductWorker(ctx context.Context)
}

type productConsumerWorker struct {
	conn     *amqp.Connection
	esClient *elasticsearch.Client
	logger   *log.Logger
	cfg      *config.Config
}

func NewProductConsumerWorker(conn *amqp.Connection, esClient *elasticsearch.Client, cfg *config.Config, logger *log.Logger) ProductConsumerWorkerInterface {
	return &productConsumerWorker{
		conn:     conn,
		esClient: esClient,
		cfg:      cfg,
		logger:   logger,
	}
}

// StartUpdateProductWorker implements [ProductConsumerWorkerInterface].
func (p *productConsumerWorker) StartUpdateProductWorker(ctx context.Context) {
	ch, err := p.conn.Channel()
	if err != nil {
		p.logger.Errorf("[ProductConsumer] StartUpdateProductWorker: failed to open a channel: %v", err)
	}

	defer ch.Close()

	productUpdate := p.cfg.PublisherName.ProductUpdate

	queue, err := ch.QueueDeclare(productUpdate, true, false, false, false, nil)
	if err != nil {
		p.logger.Errorf("[ProductConsumer] StartUpdateProductWorker: failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		p.logger.Errorf("[ProductConsumer] StartUpdateProductWorker: failed to register consumer: %v", err)
	}

	p.logger.Infof("[ProductConsumer] StartUpdateProductWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				p.logger.Infof("[ProductConsumer] StartUpdateProductWorker: message channel closed")
				return
			}

			var product entity.ProductEntity

			err := json.Unmarshal(d.Body, &product)
			if err != nil {
				p.logger.Errorf("[ProductConsumer] StartUpdateProductWorker: error decoding message: %v", err)
				continue
			}

			productJson, err := json.Marshal(&product)
			if err != nil {
				p.logger.Errorf("[ProductConsumer] StartUpdateProductWorker: error encoding product to json: %v", err)
				continue
			}

			if _, err := p.esClient.Index(
				"products",
				bytes.NewReader(productJson),
				p.esClient.Index.WithDocumentID(fmt.Sprintf("%d", product.ID)),
				p.esClient.Index.WithContext(ctx),
				p.esClient.Index.WithRefresh("true"),
			); err != nil {
				p.logger.Errorf("[ProductConsumer] StartUpdateProductWorker: error update document to elasticsearch: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// body, _ := io.ReadAll(res.Body)
			// defer res.Body.Close()

			p.logger.Infof("[ProductConsumer] StartUpdateProductWorker: product %d successfully updated to elasticsearch", product.ID)
		}
	}
}

// StartCreateProductWorker implements [ProductConsumerWorkerInterface].
func (p *productConsumerWorker) StartCreateProductWorker(ctx context.Context) {
	ch, err := p.conn.Channel()
	if err != nil {
		p.logger.Errorf("[ProductConsumer] StartCreateProductWorker: failed to open a channel: %v", err)
	}

	defer ch.Close()

	ProductCreate := p.cfg.PublisherName.ProductCreate

	queue, err := ch.QueueDeclare(ProductCreate, true, false, false, false, nil)
	if err != nil {
		p.logger.Errorf("[ProductConsumer] StartCreateProductWorker: failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		p.logger.Errorf("[ProductConsumer] StartCreateProductWorker: failed to register consumer: %v", err)
	}

	p.logger.Infof("[ProductConsumer] StartCreateProductWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				p.logger.Infof("[ProductConsumer] StartCreateProductWorker: message channel closed")
				return
			}

			var product entity.ProductEntity

			err := json.Unmarshal(d.Body, &product)
			if err != nil {
				p.logger.Errorf("[ProductConsumer] StartCreateProductWorker: error decoding message: %v", err)
				continue
			}

			productJson, err := json.Marshal(product)
			if err != nil {
				p.logger.Errorf("[ProductConsumer] StartCreateProductWorker: error encoding product to json: %v", err)
				continue
			}

			if _, err := p.esClient.Index(
				"products",
				bytes.NewReader(productJson),
				p.esClient.Index.WithDocumentID(fmt.Sprintf("%d", product.ID)),
				p.esClient.Index.WithContext(ctx),
				p.esClient.Index.WithRefresh("true"),
			); err != nil {
				p.logger.Errorf("[ProductConsumer] StartCreateProductWorker: error indexing to elasticsearch: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// body, _ := io.ReadAll(res.Body)
			// defer res.Body.Close()

			p.logger.Infof("[ProductConsumer] StartCreateProductWorker: product %d successfully indexed to elasticsearch", product.ID)
		}
	}
}

// StartDeleteProductWorker implements [ProductConsumerWorkerInterface].
func (p *productConsumerWorker) StartDeleteProductWorker(ctx context.Context) {
	ch, err := p.conn.Channel()
	if err != nil {
		p.logger.Errorf("[ProductConsumer] StartDeleteProductWorker: failed to open a channel: %v", err)
	}

	defer ch.Close()

	ProductDelete := p.cfg.PublisherName.ProductDelete

	queue, err := ch.QueueDeclare(ProductDelete, true, false, false, false, nil)
	if err != nil {
		p.logger.Errorf("[ProductConsumer] StartDeleteProductWorker: failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		p.logger.Errorf("[ProductConsumer] StartDeleteProductWorker: failed to register consumer: %v", err)
	}

	p.logger.Infof("[ProductConsumer] StartDeleteProductWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				p.logger.Infof("[ProductConsumer] StartDeleteProductWorker: message channel closed")
				return
			}

			var product entity.ProductEntity

			err := json.Unmarshal(d.Body, &product)
			if err != nil {
				p.logger.Errorf("[ProductConsumer] StartDeleteProductWorker: error decoding message: %v", err)
				continue
			}

			if _, err := p.esClient.Delete("products", fmt.Sprintf("%d", product.ID), p.esClient.Delete.WithContext(ctx)); err != nil {
				p.logger.Errorf("[ProductConsumer] StartDeleteProductWorker: error deleting from elasticsearch: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// body, _ := io.ReadAll(res.Body)
			// defer res.Body.Close()

			p.logger.Infof("[ProductConsumer] StartDeleteProductWorker: product %d successfully deleted from elasticsearch", product.ID)
		}
	}
}
