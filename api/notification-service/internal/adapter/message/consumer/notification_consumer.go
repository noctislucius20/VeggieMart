package consumer

import (
	"context"
	"encoding/json"
	"notification-service/internal/adapter/message"
	"notification-service/internal/adapter/repository"
	"notification-service/internal/core/domain/entity"
	"notification-service/internal/core/service"
	"notification-service/internal/core/service/transaction"
	"notification-service/utils"
	"notification-service/utils/ws"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type NotificationConsumerWorkerInterface interface {
	StartCreateNotificationWorker(ctx context.Context, queueName string)

	sendNotification(ctx context.Context, notification entity.NotificationEntity)
}

type notificationConsumerWorker struct {
	conn                *amqp091.Connection
	emailService        message.EmailMessageInterface
	repoNotification    repository.NotificationRepositoryInterface
	txManager           transaction.TransactionManager
	serviceNotification service.NotificationServiceInterface
	db                  *gorm.DB
	redisClient         *redis.Client
	logger              *log.Logger
}

func NewNotificationConsumerWorker(emailService message.EmailMessageInterface, repoNotification repository.NotificationRepositoryInterface, txManager transaction.TransactionManager, serviceNotification service.NotificationServiceInterface, conn *amqp091.Connection, db *gorm.DB, redisClient *redis.Client, logger *log.Logger) NotificationConsumerWorkerInterface {
	return &notificationConsumerWorker{
		emailService:        emailService,
		repoNotification:    repoNotification,
		txManager:           txManager,
		serviceNotification: serviceNotification,
		conn:                conn,
		db:                  db,
		redisClient:         redisClient,
		logger:              logger,
	}
}

// sendNotification implements [NotificationConsumerWorkerInterface].
func (n *notificationConsumerWorker) sendNotification(ctx context.Context, notification entity.NotificationEntity) {
	switch notification.NotificationType {
	case "EMAIL":
		if err := n.emailService.SendEmailNotification(*notification.ReceiverEmail, *notification.Subject, notification.Message); err != nil {
			n.logger.Errorf("[NotificationConsumer-1] sendNotification: %v", err)
			return
		}
	case "PUSH":
		if notification.ReceiverID == nil {
			n.logger.Errorf("[NotificationConsumer-2] sendNotification: receiver_id is nil")
			return
		}

		wsMsg := entity.WsRedisEntity{
			ID:         notification.ID,
			ReceiverID: *notification.ReceiverID,
			Type:       notification.NotificationType,
			Subject:    *notification.Subject,
			Message:    notification.Message,
		}

		if err := ws.PublishWebSocketMessage(ctx, n.redisClient, wsMsg, n.logger); err != nil {
			n.logger.Errorf("[NotificationConsumer-3] sendNotification: failed to publish to Redis: %v", err)
			return
		}
	default:
		n.logger.Errorf("[NotificationConsumer-4] sendNotification: %v", utils.INVALID_NOTIFICATION_TYPE)
		return
	}
}

// StartCreateNotificationWorker implements [NotificationConsumerWorkerInterface].
func (n *notificationConsumerWorker) StartCreateNotificationWorker(ctx context.Context, queueName string) {
	ch, err := n.conn.Channel()
	if err != nil {
		n.logger.Errorf("[NotificationConsumer-1] StartCreateNotificationWorker: %v", err)
		return
	}

	defer ch.Close()

	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		n.logger.Errorf("[NotificationConsumer-2] StartCreateNotificationWorker: %v", err)
		return
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		n.logger.Errorf("[NotificationConsumer-3] StartCreateNotificationWorker: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				n.logger.Infof("[NotificationConsumer-4] StartCreateNotificationWorker: %v", "channel closed")
				continue
			}

			var notification entity.NotificationEntity

			err := json.Unmarshal(d.Body, &notification)
			if err != nil {
				n.logger.Errorf("[NotificationConsumer-5] StartCreateNotificationWorker: %v", err)
				continue
			}

			notification.Status = "PENDING"
			if notification.NotificationType == "EMAIL" {
				notification.Status = "SENT"
			}

			if err := n.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
				notificationId, err := n.repoNotification.CreateNotification(txCtx, notification)
				if err != nil {
					return err
				}

				notification.ID = notificationId

				return nil
			}); err != nil {
				continue
			}

			go n.sendNotification(ctx, notification)

			// body, _ := n.ReadAll(res.Body)
			// defer res.Body.Close()

			n.logger.Infof("[NotificationConsumer-8] StartCreateNotificationWorker: email has been sent to %v", *notification.ReceiverEmail)
		}
	}
}
