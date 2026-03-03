package publisher

import (
	"context"
	"errors"
	"fmt"
	"product-service/internal/adapter/repository"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/service/transaction"
	"product-service/utils"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type StartPublisherWorkerInterface interface {
	StartPublisherWorker(ctx context.Context)
}

type startPublisherWorker struct {
	db         *gorm.DB
	repoOutbox repository.OutboxEventInterface
	txManager  transaction.TransactionManager
	conn       *amqp.Connection
	logger     *log.Logger
}

// StartPublisherWorker implements StartPublisherWorkerInterface.
func (s *startPublisherWorker) StartPublisherWorker(ctx context.Context) {
	jobChan := make(chan entity.OutboxEventEntity, 100)

	var wg sync.WaitGroup

	wg.Go(func() {
		s.startPoller(ctx, jobChan)
	})

	workerCount := 5
	for range workerCount {
		wg.Go(func() {
			s.startPublisher(ctx, jobChan)
		})
	}

	close(jobChan)

	wg.Wait()

}

func (s *startPublisherWorker) startPoller(ctx context.Context, jobs chan<- entity.OutboxEventEntity) {
	idleDelay := 2 * time.Second
	busyDelay := 20 * time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(busyDelay):
			var outboxes, err = []entity.OutboxEventEntity{}, errors.New("")

			if err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
				outboxes, err = s.repoOutbox.GetAllPendingEvent(txCtx)
				if err != nil {
					return err
				}

				return nil
			}); err != nil {
				s.logger.Errorf("[StartPublisherWorker-2] startPoller: %v", err.Error())
				return
			}

			if len(outboxes) == 0 {
				time.Sleep(idleDelay)
				continue
			}

			for _, outbox := range outboxes {
				select {
				case jobs <- outbox:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

func (s *startPublisherWorker) startPublisher(ctx context.Context, jobs <-chan entity.OutboxEventEntity) {
	ch, err := s.conn.Channel()
	if err != nil {
		s.logger.Errorf("[StartPublisherWorker-1] startPublisher: %v", err.Error())
		return
	}

	defer ch.Close()

	if err := ch.Confirm(false); err != nil {
		s.logger.Errorf("[StartPublisherWorker-2] startPublisher: %v", err.Error())
		return
	}

	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	for {
		select {
		case <-ctx.Done():
			return
		case outbox, ok := <-jobs:
			if !ok {
				s.logger.Infof("[StartPublisherWorker-3] startPublisher: job channel closed")
				return
			}

			if _, err = ch.QueueDeclare(outbox.EventType, true, false, false, false, nil); err != nil {
				s.logger.Errorf("[StartPublisherWorker-4] startPublisher: %v", err.Error())
				return
			}

			if err := s.publishOne(ctx, ch, confirms, outbox); err != nil {
				if err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
					if err := s.repoOutbox.UpdateFailedEvent(txCtx, []int64{outbox.ID}); err != nil {
						return err
					}

					return nil
				}); err != nil {
					s.logger.Errorf("[StartPublisherWorker-5] startPublisher: %v", err.Error())
					return
				}

				continue
			}

			if err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
				if err := s.repoOutbox.UpdatePublishedEvent(txCtx, []int64{outbox.ID}); err != nil {
					return err
				}

				return nil
			}); err != nil {
				s.logger.Errorf("[StartPublisherWorker-6] startPublisher: %v", err.Error())
				return
			}
		}
	}
}

func (s *startPublisherWorker) publishOne(ctx context.Context, ch *amqp.Channel, confirms <-chan amqp.Confirmation, outbox entity.OutboxEventEntity) error {
	if err := ch.PublishWithContext(
		ctx,
		"",
		outbox.EventType,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(outbox.Payload),
			MessageId:   fmt.Sprintf("%d", outbox.ID),
		}); err != nil {
		s.logger.Errorf("[StartPublisherWorker-1] publishOne: %v", err.Error())
		return err
	}

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		err := errors.New(utils.SERVICE_UNAVAILABLE)
		s.logger.Errorf("[StartPublisherWorker-2] publishOne: %v", err.Error())

		return err
	case confirm := <-confirms:
		if !confirm.Ack {
			s.logger.Errorf("[StartPublisherWorker-3] publishOne: publish id %d failed", outbox.ID)

			return errors.New(utils.SERVICE_UNAVAILABLE)
		}
	case <-timer.C:
		s.logger.Errorf("[StartPublisherWorker-4] publishOne: publish id %d timeout", outbox.ID)

		return errors.New(utils.TIMEOUT_LIMIT_EXCEEDED)
	}

	return nil
}

func NewStartPublisherWorker(db *gorm.DB, conn *amqp.Connection, repoOutbox repository.OutboxEventInterface, logger *log.Logger) StartPublisherWorkerInterface {
	return &startPublisherWorker{db: db, conn: conn, repoOutbox: repoOutbox, logger: logger}
}
