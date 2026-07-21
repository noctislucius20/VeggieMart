package service

import (
	"context"
	"notification-service/internal/adapter/repository"
	"notification-service/internal/core/domain/entity"
	"notification-service/internal/core/service/transaction"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type NotificationServiceInterface interface {
	GetAllNotifications(ctx context.Context, query entity.NotificationQueryString) ([]entity.NotificationEntity, int64, int64, error)
	GetAllPushNotification(ctx context.Context, query entity.NotificationQueryString) ([]entity.NotificationEntity, int64, int64, error)
	GetNotificationById(ctx context.Context, notificationId int64) (*entity.NotificationEntity, error)
	MarkAsSentNotification(ctx context.Context, notificationId int64) error
	MarkAsReadNotification(ctx context.Context, notificationId int64) error
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

// GetAllPushNotification implements [NotificationServiceInterface].
func (n *notificationService) GetAllPushNotification(ctx context.Context, query entity.NotificationQueryString) ([]entity.NotificationEntity, int64, int64, error) {
	var (
		notifications []entity.NotificationEntity
		countData     int64
		totalPages    int64
	)

	if err := n.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		notificationEntities, count, pages, err := n.repo.GetAllPushNotification(txCtx, query)
		if err != nil {
			return nil
		}

		if len(notificationEntities) == 0 {
			return nil
		}

		if err := n.repo.MarkAllAsSentNotification(txCtx, query.UserID); err != nil {
			return err
		}

		notifications, countData, totalPages = notificationEntities, count, pages

		return nil
	}); err != nil {
		n.logger.Errorf("[NotificationService-1] GetAllPushNotification: %v", err)
		return nil, 0, 0, err
	}

	return notifications, countData, totalPages, nil
}

// MarkAsReadNotification implements [NotificationServiceInterface].
func (n *notificationService) MarkAsReadNotification(ctx context.Context, notificationId int64) error {
	if err := n.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if _, err := n.repo.GetNotificationById(txCtx, notificationId); err != nil {
			return err
		}

		if err := n.repo.MarkAsReadNotification(txCtx, notificationId); err != nil {
			return err
		}

		return nil
	}); err != nil {
		n.logger.Errorf("[NotificationService-1] MarkAsReadNotification: %v", err)
		return err
	}

	return nil
}

// MarkAsSentNotification implements [NotificationServiceInterface].
func (n *notificationService) MarkAsSentNotification(ctx context.Context, notificationId int64) error {
	if err := n.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := n.repo.MarkAsSentNotification(txCtx, notificationId); err != nil {
			return err
		}

		return nil
	}); err != nil {
		n.logger.Errorf("[NotificationService-1] MarkAsSentNotification: %v", err)
		return err
	}

	return nil
}

// GetNotificationById implements [NotificationServiceInterface].
func (n *notificationService) GetNotificationById(ctx context.Context, notificationId int64) (*entity.NotificationEntity, error) {
	notification := &entity.NotificationEntity{}

	if err := n.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		notificationEntity, err := n.repo.GetNotificationById(txCtx, notificationId)
		if err != nil {
			return err
		}

		notification = notificationEntity

		return nil
	}); err != nil {
		n.logger.Errorf("[NotificationService-1] GetNotificationById: %v", err)
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
		if err := n.repo.MarkAllAsSentNotification(txCtx, query.UserID); err != nil {
			return err
		}

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
		n.logger.Errorf("[NotificationService-1] GetAllNotifications: %v", err)
		return nil, 0, 0, err
	}

	return notifications, countData, totalPages, nil
}
