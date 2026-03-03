package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"product-service/config"
	"product-service/internal/core/domain/entity"

	"github.com/labstack/gommon/log"
)

type ProductConsumerWorkerInterface interface {
	StartCreateProductWorker(ctx context.Context)
	StartDeleteProductWorker(ctx context.Context)
	StartUpdateProductWorker(ctx context.Context)
}

type productConsumerWorker struct {
	cfg *config.Config
}

// StartUpdateProductWorker implements [ProductConsumerWorkerInterface].
func (p *productConsumerWorker) StartUpdateProductWorker(ctx context.Context) {
	conn, err := p.cfg.NewRabbitMQ()
	if err != nil {
		log.Errorf("[ProductConsumer-1] StartUpdateProductWorker: failed to connect to RabbitMQ: %v", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[ProductConsumer-2] StartUpdateProductWorker: failed to open a channel: %v", err)
	}

	defer ch.Close()

	productUpdate := config.NewConfig().PublisherName.ProductUpdate

	queue, err := ch.QueueDeclare(productUpdate, true, false, false, false, nil)
	if err != nil {
		log.Errorf("[ProductConsumer-3] StartUpdateProductWorker: failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Errorf("[ProductConsumer-4] StartUpdateProductWorker: failed to register consumer: %v", err)
	}

	esClient, err := config.NewConfig().NewElasticSearchClient()
	if err != nil {
		log.Errorf("[ProductConsumer-5] StartUpdateProductWorker: failed to initialize elasticsearch client: %v", err)
	}

	log.Infof("[ProductConsumer-6] StartUpdateProductWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				log.Infof("[ProductConsumer-7] StartUpdateProductWorker: message channel closed")
				return
			}

			var product entity.ProductEntity

			err := json.Unmarshal(d.Body, &product)
			if err != nil {
				log.Errorf("[ProductConsumer-8] StartUpdateProductWorker: error decoding message: %v", err)
				continue
			}

			productJson, err := json.Marshal(&product)
			if err != nil {
				log.Errorf("[ProductConsumer-9] StartUpdateProductWorker: error encoding product to json: %v", err)
				continue
			}

			if _, err := esClient.Index(
				"products",
				bytes.NewReader(productJson),
				esClient.Index.WithDocumentID(fmt.Sprintf("%d", product.ID)),
				esClient.Index.WithContext(ctx),
				esClient.Index.WithRefresh("true"),
			); err != nil {
				log.Errorf("[ProductConsumer-10] StartUpdateProductWorker: error update document to elasticsearch: %v", err)
				continue
			}

			// body, _ := io.ReadAll(res.Body)
			// defer res.Body.Close()

			log.Infof("[ProductConsumer-11] StartUpdateProductWorker: product %d successfully updated to elasticsearch", product.ID)
		}
	}
}

// StartCreateProductWorker implements [ProductConsumerWorkerInterface].
func (p *productConsumerWorker) StartCreateProductWorker(ctx context.Context) {
	conn, err := p.cfg.NewRabbitMQ()
	if err != nil {
		log.Errorf("[ProductConsumer-1] StartCreateProductWorker: failed to connect to RabbitMQ: %v", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[ProductConsumer-2] StartCreateProductWorker: failed to open a channel: %v", err)
	}

	defer ch.Close()

	ProductCreate := config.NewConfig().PublisherName.ProductCreate

	queue, err := ch.QueueDeclare(ProductCreate, true, false, false, false, nil)
	if err != nil {
		log.Errorf("[ProductConsumer-3] StartCreateProductWorker: failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Errorf("[ProductConsumer-4] StartCreateProductWorker: failed to register consumer: %v", err)
	}

	esClient, err := config.NewConfig().NewElasticSearchClient()
	if err != nil {
		log.Errorf("[ProductConsumer-5] StartCreateProductWorker: failed to initialize elasticsearch client: %v", err)
	}

	log.Infof("[ProductConsumer-6] StartCreateProductWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				log.Infof("[ProductConsumer-7] StartCreateProductWorker: message channel closed")
				return
			}

			var product entity.ProductEntity

			err := json.Unmarshal(d.Body, &product)
			if err != nil {
				log.Errorf("[ProductConsumer-8] StartCreateProductWorker: error decoding message: %v", err)
				continue
			}

			productJson, err := json.Marshal(product)
			if err != nil {
				log.Errorf("[ProductConsumer-9] StartCreateProductWorker: error encoding product to json: %v", err)
				continue
			}

			if _, err := esClient.Index(
				"products",
				bytes.NewReader(productJson),
				esClient.Index.WithDocumentID(fmt.Sprintf("%d", product.ID)),
				esClient.Index.WithContext(ctx),
				esClient.Index.WithRefresh("true"),
			); err != nil {
				log.Errorf("[ProductConsumer-10] StartCreateProductWorker: error indexing to elasticsearch: %v", err)
				continue
			}

			// body, _ := io.ReadAll(res.Body)
			// defer res.Body.Close()

			log.Infof("[ProductConsumer-11] StartCreateProductWorker: product %d successfully indexed to elasticsearch", product.ID)
		}
	}
}

// StartDeleteProductWorker implements [ProductConsumerWorkerInterface].
func (p *productConsumerWorker) StartDeleteProductWorker(ctx context.Context) {
	conn, err := p.cfg.NewRabbitMQ()
	if err != nil {
		log.Errorf("[ProductConsumer-1] StartDeleteProductWorker: failed to connect to RabbitMQ: %v", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[ProductConsumer-2] StartDeleteProductWorker: failed to open a channel: %v", err)
	}

	defer ch.Close()

	ProductDelete := config.NewConfig().PublisherName.ProductDelete

	queue, err := ch.QueueDeclare(ProductDelete, true, false, false, false, nil)
	if err != nil {
		log.Errorf("[ProductConsumer-3] StartDeleteProductWorker: failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Errorf("[ProductConsumer-4] StartDeleteProductWorker: failed to register consumer: %v", err)
	}

	esClient, err := config.NewConfig().NewElasticSearchClient()
	if err != nil {
		log.Errorf("[ProductConsumer-5] StartDeleteProductWorker: failed to initialize elasticsearch client: %v", err)
	}

	log.Infof("[ProductConsumer-6] StartDeleteProductWorker: waiting for messages. to exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				log.Infof("[ProductConsumer-7] StartDeleteProductWorker: message channel closed")
				return
			}

			var product entity.ProductEntity

			err := json.Unmarshal(d.Body, &product)
			if err != nil {
				log.Errorf("[ProductConsumer-8] StartDeleteProductWorker: error decoding message: %v", err)
				continue
			}

			if _, err := esClient.Delete("products", fmt.Sprintf("%d", product.ID), esClient.Delete.WithContext(ctx)); err != nil {
				log.Errorf("[ProductConsumer-9] StartDeleteProductWorker: error deleting from elasticsearch: %v", err)
				continue
			}

			// body, _ := io.ReadAll(res.Body)
			// defer res.Body.Close()

			log.Infof("[ProductConsumer-10] StartDeleteProductWorker: product %d successfully deleted from elasticsearch", product.ID)
		}
	}
}

func NewProductConsumerWorker(cfg *config.Config) ProductConsumerWorkerInterface {
	return &productConsumerWorker{cfg: cfg}
}
