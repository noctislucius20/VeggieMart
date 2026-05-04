package service

import (
	"context"
	"fmt"
	"notification-service/internal/adapter/repository"
	"notification-service/internal/core/domain/entity"
	"notification-service/internal/core/service/transaction"
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
	repo      repository.NotificationRepositoryInterface
	db        *gorm.DB
	txManager transaction.TransactionManager
	logger    *log.Logger
}

func NewNotificationService(repo repository.NotificationRepositoryInterface, txManager transaction.TransactionManager, db *gorm.DB, logger *log.Logger) NotificationServiceInterface {
	return &notificationService{
		repo:      repo,
		txManager: txManager,
		db:        db,
		logger:    logger,
	}
}

// TODO add cache notification

// MarkAsReadNotification implements [NotificationServiceInterface].
func (n *notificationService) MarkAsReadNotification(ctx context.Context, notificationId uint) error {
	if err := n.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if _, err := n.repo.GetNotificationById(txCtx, notificationId); err != nil {
			return err
		}

		if err := n.repo.MarkAsReadNotification(txCtx, notificationId); err != nil {
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

	if err := n.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
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

		if _, err := n.repo.GetNotificationById(txCtx, notification.ID); err != nil {
			return err
		}

		if err := n.repo.MarkAsSentNotification(txCtx, notification.ID); err != nil {
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

	if err := n.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		notificationEntity, err := n.repo.GetNotificationById(txCtx, notificationId)
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

	if err := n.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		notificationEntities, count, pages, err := n.repo.GetAllNotifications(txCtx, query)
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
