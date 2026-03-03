package consumer

import (
	"context"
	"encoding/json"
	"notification-service/internal/adapter/message"
	"notification-service/internal/adapter/repository"
	"notification-service/internal/core/domain/entity"
	"notification-service/internal/core/service"
	"notification-service/utils"

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
	serviceNotification service.NotificationServiceInterface
	db                  *gorm.DB
	logger              *log.Logger
}

// sendNotification implements [NotificationConsumerWorkerInterface].
func (n *notificationConsumerWorker) sendNotification(ctx context.Context, notification entity.NotificationEntity) {
	switch notification.NotificationType {
	case "EMAIL":
		if err := n.emailService.SendEmailNotification(*notification.ReceiverEmail, *notification.Subject, notification.Message); err != nil {
			n.logger.Errorf("[NotificationConsumer-1] sendNotification: %v", err.Error())
			return
		}
	case "PUSH":
		n.serviceNotification.SendPushNotification(ctx, notification)
		return
	default:
		n.logger.Errorf("[NotificationConsumer-2] sendNotification: %v", utils.INVALID_NOTIFICATION_TYPE)
		return
	}
}

// StartCreateNotificationWorker implements [NotificationConsumerWorkerInterface].
func (n *notificationConsumerWorker) StartCreateNotificationWorker(ctx context.Context, queueName string) {
	ch, err := n.conn.Channel()
	if err != nil {
		n.logger.Errorf("[NotificationConsumer-1] StartCreateNotificationWorker: %v", err.Error())
		return
	}

	defer ch.Close()

	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		n.logger.Errorf("[NotificationConsumer-2] StartCreateNotificationWorker: %v", err.Error())
		return
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		n.logger.Errorf("[NotificationConsumer-3] StartCreateNotificationWorker: %v", err.Error())
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
				n.logger.Errorf("[NotificationConsumer-5] StartCreateNotificationWorker: %v", err.Error())
				continue
			}

			notification.Status = "PENDING"
			if notification.NotificationType == "EMAIL" {
				notification.Status = "SENT"
			}

			if err := n.db.Transaction(func(tx *gorm.DB) error {
				if err := n.repoNotification.CreateNotification(ctx, notification, tx); err != nil {
					return err
				}

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

func NewNotificationConsumerWorker(emailService message.EmailMessageInterface, repoNotification repository.NotificationRepositoryInterface, serviceNotification service.NotificationServiceInterface, conn *amqp091.Connection, db *gorm.DB, logger *log.Logger) NotificationConsumerWorkerInterface {
	return &notificationConsumerWorker{
		emailService:        emailService,
		repoNotification:    repoNotification,
		serviceNotification: serviceNotification,
		conn:                conn,
		db:                  db,
		logger:              logger,
	}
}
