package service

import (
	"context"
	"fmt"
	"notification-service/internal/adapter/repository"
	"notification-service/internal/core/domain/entity"
	"notification-service/utils"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type NotificationServiceInterface interface {
	GetAllNotifications(ctx context.Context, query entity.NotificationQueryString) ([]entity.NotificationEntity, int64, int64, error)
	GetNotificationById(ctx context.Context, notificationId uint) (*entity.NotificationEntity, error)
	SendPushNotification(ctx context.Context, notification entity.NotificationEntity)
	MarkAsReadNotification(ctx context.Context, notificationId uint) error
}

type notificationService struct {
	repo   repository.NotificationRepositoryInterface
	db     *gorm.DB
	logger *log.Logger
}

// MarkAsReadNotification implements [NotificationServiceInterface].
func (n *notificationService) MarkAsReadNotification(ctx context.Context, notificationId uint) error {
	if err := n.db.Transaction(func(tx *gorm.DB) error {
		if _, err := n.repo.GetNotificationById(ctx, notificationId, tx); err != nil {
			return err
		}

		if err := n.repo.MarkAsReadNotification(ctx, notificationId, tx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		n.logger.Errorf("[NotificationService-1] MarkAsReadNotification: %v", err.Error())
		return err
	}

	return nil
}

// SendPushNotification implements [NotificationServiceInterface].
func (n *notificationService) SendPushNotification(ctx context.Context, notification entity.NotificationEntity) {
	if notification.ReceiverID == nil {
		return
	}

	if err := n.db.Transaction(func(tx *gorm.DB) error {
		conn := utils.GetWebSocketConn(*notification.ReceiverID)
		if conn == nil {
			err := fmt.Errorf("%v, ID = %d", utils.DATA_NOT_FOUND, *notification.ReceiverID)
			return err
		}

		msg := map[string]any{
			"type":    notification.NotificationType,
			"subject": notification.Subject,
			"message": notification.Message,
			"sent_at": notification.SentAt,
		}

		if err := conn.WriteJSON(msg); err != nil {
			return err
		}

		if _, err := n.repo.GetNotificationById(ctx, notification.ID, tx); err != nil {
			return err
		}

		if err := n.repo.MarkAsSentNotification(ctx, notification.ID, tx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		n.logger.Errorf("[NotificationService-2] SendPushNotification: %v", err.Error())
		return
	}
}

// GetNotificationById implements [NotificationServiceInterface].
func (n *notificationService) GetNotificationById(ctx context.Context, notificationId uint) (*entity.NotificationEntity, error) {
	notification := &entity.NotificationEntity{}

	if err := n.db.Transaction(func(tx *gorm.DB) error {
		notificationEntity, err := n.repo.GetNotificationById(ctx, notificationId, tx)
		if err != nil {
			return err
		}

		notification = notificationEntity

		return nil
	}); err != nil {
		n.logger.Errorf("[NotificationService-1] GetNotificationById: %v", err.Error())
		return nil, err
	}

	return notification, nil
}

// GetAllNotifications implements [NotificationServiceInterface].
func (n *notificationService) GetAllNotifications(ctx context.Context, query entity.NotificationQueryString) ([]entity.NotificationEntity, int64, int64, error) {
	var (
		notifications []entity.NotificationEntity
		countData     int64
		totalPages    int64
	)

	if err := n.db.Transaction(func(tx *gorm.DB) error {
		notificationEntities, count, pages, err := n.repo.GetAllNotifications(ctx, query, tx)
		if err != nil {
			return nil
		}

		if len(notificationEntities) == 0 {
			return nil
		}

		notifications, countData, totalPages = notificationEntities, count, pages

		return nil
	}); err != nil {
		n.logger.Errorf("[NotificationService-1] GetAllNotifications: %v", err.Error())
		return nil, 0, 0, err
	}

	return notifications, countData, totalPages, nil
}

func NewNotificationService(repo repository.NotificationRepositoryInterface, db *gorm.DB, logger *log.Logger) NotificationServiceInterface {
	return &notificationService{repo: repo, db: db, logger: logger}
}
