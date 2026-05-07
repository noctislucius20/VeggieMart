package publisher

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"user-service/internal/adapter/repository"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/service/transaction"
	"user-service/utils"

	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type StartPublisherWorkerInterface interface {
	StartPublisherWorker(ctx context.Context)
}

type startPublisherWorker struct {
	repoOutbox repository.OutboxEventInterface
	txManager  transaction.TransactionManager
	conn       *amqp.Connection
	logger     *log.Logger
}

func NewStartPublisherWorker(conn *amqp.Connection, txManager transaction.TransactionManager, repoOutbox repository.OutboxEventInterface, logger *log.Logger) StartPublisherWorkerInterface {
	return &startPublisherWorker{
		conn:       conn,
		txManager:  txManager,
		repoOutbox: repoOutbox,
		logger:     logger,
	}
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

	wg.Wait()

	close(jobChan)
}

func (s *startPublisherWorker) startPoller(ctx context.Context, jobs chan<- entity.OutboxEventEntity) {
	idleDelay := 2 * time.Second
	busyDelay := 20 * time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(busyDelay):
			var (
				outboxes []entity.OutboxEventEntity
			)

			if err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
				outboxEntities, err := s.repoOutbox.GetAllPendingEvent(txCtx)
				if err != nil {
					return err
				}

				outboxes = outboxEntities

				return nil
			}); err != nil {
				s.logger.Errorf("[StartPublisherWorker-1] startPoller: %v", err)
				time.Sleep(idleDelay)
				continue
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
		s.logger.Fatalf("[StartPublisherWorker-1] startPublisher: %v", err)
		return
	}

	defer ch.Close()

	if err := ch.Confirm(false); err != nil {
		s.logger.Fatalf("[StartPublisherWorker-2] startPublisher: %v", err)
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
				continue
			}

			if _, err = ch.QueueDeclare(outbox.EventType, true, false, false, false, nil); err != nil {
				s.logger.Errorf("[StartPublisherWorker-4] startPublisher: %v", err)
				continue
			}

			if err := s.publishOne(ctx, ch, confirms, outbox); err != nil {
				if err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
					if err := s.repoOutbox.UpdateFailedEvent(txCtx, []int64{outbox.ID}); err != nil {
						return err
					}

					return nil
				}); err != nil {
					s.logger.Errorf("[StartPublisherWorker-5] startPublisher: %v", err)
					continue
				}

				continue
			}

			if err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
				if err := s.repoOutbox.UpdatePublishedEvent(txCtx, []int64{outbox.ID}); err != nil {
					return err
				}

				return nil
			}); err != nil {
				s.logger.Errorf("[StartPublisherWorker-6] startPublisher: %v", err)
				continue
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
		s.logger.Errorf("[StartPublisherWorker-1] publishOne: %v", err)
		return err
	}

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		err := errors.New(utils.SERVICE_UNAVAILABLE)
		s.logger.Errorf("[StartPublisherWorker-2] publishOne: %v", err)

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
