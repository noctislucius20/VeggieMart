package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"notification-service/internal/core/domain/entity"
	"notification-service/internal/core/domain/model"
	"notification-service/utils"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type NotificationRepositoryInterface interface {
	GetAllNotifications(ctx context.Context, query entity.NotificationQueryString) ([]entity.NotificationEntity, int64, int64, error)
	GetAllPushNotification(ctx context.Context, query entity.NotificationQueryString) ([]entity.NotificationEntity, int64, int64, error)
	GetNotificationById(ctx context.Context, notificationId int64) (*entity.NotificationEntity, error)
	CreateNotification(ctx context.Context, notification entity.NotificationEntity) (int64, error)
	MarkAsSentNotification(ctx context.Context, notificationId int64) error
	MarkAsReadNotification(ctx context.Context, notificationId int64) error
	MarkAllAsSentNotification(ctx context.Context, userId int64) error

	getDB(ctx context.Context) *gorm.DB
}

type notificationRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewNotificationRepository(db *gorm.DB, logger *log.Logger) NotificationRepositoryInterface {
	return &notificationRepository{db: db, logger: logger}
}

// getDB implements [NotificationRepositoryInterface].
func (n *notificationRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}

	return n.db
}

// GetAllPushNotification implements [NotificationRepositoryInterface].
func (n *notificationRepository) GetAllPushNotification(ctx context.Context, query entity.NotificationQueryString) ([]entity.NotificationEntity, int64, int64, error) {
	var (
		db                 = n.getDB(ctx)
		modelNotifications []model.Notification
	)

	var countData int64

	offset := (query.Page - 1) * query.Limit

	sqlMain := db.WithContext(ctx).
		Omit("receiver_email", "created_at", "updated_at", "deleted_at").
		Where("subject ILIKE ? OR message ILIKE ? OR status ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%", "%"+query.Status+"%")

	if query.UserID > 0 {
		sqlMain = sqlMain.
			Where("receiver_id = ?", query.UserID).
			Where("notification_method = ?", "PUSH")
	} else {
		sqlMain = sqlMain.Where("notification_method = ?", "ADMIN_PUSH")
	}

	if query.IsRead {
		sqlMain = sqlMain.Where("read_at IS NOT NULL")
	}

	if err := sqlMain.Model(&modelNotifications).
		Count(&countData).Error; err != nil {
		n.logger.Errorf("[NotificationRepository] GetAllPushNotification: %v", err)
		return nil, 0, 0, err
	}

	orderSort := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)

	if err := sqlMain.Order(orderSort).
		Limit(int(query.Limit)).
		Offset(int(offset)).
		Find(&modelNotifications).Error; err != nil {
		n.logger.Errorf("[NotificationRepository] GetAllPushNotification: %v", err)
		return nil, 0, 0, err
	}

	var notifications []entity.NotificationEntity
	for _, modelNotification := range modelNotifications {
		notifications = append(notifications, entity.NotificationEntity{
			ID:                 modelNotification.ID,
			NotificationType:   modelNotification.NotificationType,
			NotificationTypeID: modelNotification.NotificationTypeID,
			Subject:            modelNotification.Subject,
			Status:             modelNotification.Status,
			SentAt:             modelNotification.SentAt,
			ReadAt:             modelNotification.ReadAt,
			Message:            modelNotification.Message,
			NotificationMethod: modelNotification.NotificationMethod,
		})
	}

	totalPage := int64(math.Ceil(float64(countData) / float64(query.Limit)))

	return notifications, countData, totalPage, nil
}

// MarkAllAsSentNotification implements [NotificationRepositoryInterface].
func (n *notificationRepository) MarkAllAsSentNotification(ctx context.Context, userId int64) error {
	var (
		db                = n.getDB(ctx)
		now               = time.Now()
		modelNotification = model.Notification{
			Status: "SENT",
			SentAt: &now,
		}
	)

	tx := db.WithContext(ctx).
		Where("status = ?", "PENDING")

	if userId > 0 {
		tx = tx.
			Where("receiver_id = ?", userId).
			Where("notification_method = ?", "PUSH")
	} else {
		tx = tx.Where("notification_method = ?", "ADMIN_PUSH")
	}

	tx = tx.Updates(&modelNotification)

	if tx.Error != nil {
		n.logger.Errorf("[NotificationRepository] MarkAllAsSentNotification: %v", tx.Error)
		return tx.Error
	}

	return nil
}

// MarkAsReadNotification implements [NotificationRepositoryInterface].
func (n *notificationRepository) MarkAsReadNotification(ctx context.Context, notificationId int64) error {
	var (
		db                = n.getDB(ctx)
		now               = time.Now()
		modelNotification = model.Notification{
			ID:     notificationId,
			ReadAt: &now,
		}
	)

	tx := db.WithContext(ctx).Updates(&modelNotification)
	if tx.Error != nil {
		n.logger.Errorf("[NotificationRepository] MarkAsReadNotification: %v", tx.Error)
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		err := utils.ErrDataNotFound
		n.logger.Errorf("[NotificationRepository] MarkAsReadNotification: %v", err)
		return err
	}

	return nil
}

// MarkAsSentNotification implements [NotificationRepositoryInterface].
func (n *notificationRepository) MarkAsSentNotification(ctx context.Context, notificationId int64) error {
	var (
		db                = n.getDB(ctx)
		now               = time.Now()
		modelNotification = model.Notification{
			ID:     notificationId,
			SentAt: &now,
			Status: "SENT",
		}
	)

	tx := db.WithContext(ctx).Updates(&modelNotification)
	if tx.Error != nil {
		n.logger.Errorf("[NotificationRepository] MarkAsSentNotification: %v", tx.Error)
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		err := utils.ErrDataNotFound
		n.logger.Errorf("[NotificationRepository] MarkAsSentNotification: %v", err)
		return err
	}

	return nil
}

// CreateNotification implements [NotificationRepositoryInterface].
func (n *notificationRepository) CreateNotification(ctx context.Context, notification entity.NotificationEntity) (int64, error) {
	var (
		db                = n.getDB(ctx)
		now               = time.Now()
		modelNotification = model.Notification{
			NotificationType:   notification.NotificationType,
			NotificationTypeID: notification.NotificationTypeID,
			ReceiverID:         notification.ReceiverID,
			Subject:            notification.Subject,
			Status:             notification.Status,
			ReadAt:             notification.ReadAt,
			Message:            notification.Message,
			NotificationMethod: notification.NotificationMethod,
		}
	)

	if notification.NotificationMethod == "EMAIL" {
		modelNotification.SentAt = &now
		modelNotification.ReceiverEmail = notification.ReceiverEmail
	}

	if err := db.WithContext(ctx).Create(&modelNotification).Error; err != nil {
		n.logger.Errorf("[NotificationRepository] CreateNotification: %v", err)
		return 0, err
	}

	return modelNotification.ID, nil
}

// GetNotificationById implements [NotificationRepositoryInterface].
func (n *notificationRepository) GetNotificationById(ctx context.Context, notificationId int64) (*entity.NotificationEntity, error) {
	var (
		db                = n.getDB(ctx)
		modelNotification model.Notification
	)

	sqlMain := db.WithContext(ctx).
		Where("id = ?", notificationId).
		Omit("created_at", "updated_at", "deleted_at")

	if err := sqlMain.First(&modelNotification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = utils.ErrDataNotFound
		}
		n.logger.Errorf("[NotificationRepository] GetNotificationById: %v", err)
		return nil, err
	}

	return &entity.NotificationEntity{
		ID:                 modelNotification.ID,
		Subject:            modelNotification.Subject,
		Status:             modelNotification.Status,
		SentAt:             modelNotification.SentAt,
		ReadAt:             modelNotification.ReadAt,
		Message:            modelNotification.Message,
		NotificationMethod: modelNotification.NotificationMethod,
	}, nil
}

// GetAllNotifications implements [NotificationRepositoryInterface].
func (n *notificationRepository) GetAllNotifications(ctx context.Context, query entity.NotificationQueryString) ([]entity.NotificationEntity, int64, int64, error) {
	var (
		db                 = n.getDB(ctx)
		modelNotifications []model.Notification
	)

	var countData int64

	offset := (query.Page - 1) * query.Limit

	sqlMain := db.WithContext(ctx).
		Select("id", "subject", "status", "sent_at").
		Where("subject ILIKE ? OR message ILIKE ? OR status ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%", "%"+query.Status+"%")

	if query.UserID != 0 {
		sqlMain = sqlMain.Where("receiver_id = ?", query.UserID)
	}

	if query.IsRead {
		sqlMain = sqlMain.Where("read_at IS NOT NULL")
	}

	if err := sqlMain.Model(&modelNotifications).Count(&countData).Error; err != nil {
		n.logger.Errorf("[NotificationRepository] GetAllNotifications: %v", err)
		return nil, 0, 0, err
	}

	orderSort := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)

	totalPages := int(math.Ceil(float64(countData) / float64(query.Limit)))
	if err := sqlMain.Order(orderSort).
		Limit(int(query.Limit)).
		Offset(int(offset)).
		Find(&modelNotifications).Error; err != nil {
		n.logger.Errorf("[NotificationRepository] GetAllNotifications: %v", err)
		return nil, 0, 0, err
	}

	entities := []entity.NotificationEntity{}
	for _, val := range modelNotifications {
		entities = append(entities, entity.NotificationEntity{
			ID:      val.ID,
			Subject: val.Subject,
			Status:  val.Status,
			SentAt:  val.SentAt,
		})
	}

	return entities, countData, int64(totalPages), nil
}
