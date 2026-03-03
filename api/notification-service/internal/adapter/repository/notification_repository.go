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
	GetAllNotifications(ctx context.Context, query entity.NotificationQueryString, db *gorm.DB) ([]entity.NotificationEntity, int64, int64, error)
	GetNotificationById(ctx context.Context, notificationId uint, db *gorm.DB) (*entity.NotificationEntity, error)
	CreateNotification(ctx context.Context, notification entity.NotificationEntity, db *gorm.DB) error
	MarkAsSentNotification(ctx context.Context, notificationId uint, db *gorm.DB) error
	MarkAsReadNotification(ctx context.Context, notificationId uint, db *gorm.DB) error
}

type notificationRepository struct {
	logger *log.Logger
}

// MarkAsReadNotification implements [NotificationRepositoryInterface].
func (n *notificationRepository) MarkAsReadNotification(ctx context.Context, notificationId uint, db *gorm.DB) error {
	now := time.Now()

	modelNotification := model.Notification{
		ID:     notificationId,
		ReadAt: &now,
	}

	if err := db.WithContext(ctx).Updates(&modelNotification).Error; err != nil {
		n.logger.Errorf("[NotificationRepository-1] MarkAsReadNotification: %v", err.Error())
		return err
	}

	return nil
}

// MarkAsSentNotification implements [NotificationRepositoryInterface].
func (n *notificationRepository) MarkAsSentNotification(ctx context.Context, notificationId uint, db *gorm.DB) error {
	modelNotification := model.Notification{
		ID:     notificationId,
		Status: "SENT",
	}

	if err := db.WithContext(ctx).Updates(&modelNotification).Error; err != nil {
		n.logger.Errorf("[NotificationRepository-1] MarkAsSentNotification: %v", err.Error())
		return err
	}

	return nil
}

// CreateNotification implements [NotificationRepositoryInterface].
func (n *notificationRepository) CreateNotification(ctx context.Context, notification entity.NotificationEntity, db *gorm.DB) error {
	now := time.Now()

	modelNotification := model.Notification{
		ReceiverID:       notification.ReceiverID,
		Subject:          notification.Subject,
		Status:           notification.Status,
		SentAt:           &now,
		ReadAt:           notification.ReadAt,
		Message:          notification.Message,
		NotificationType: notification.NotificationType,
	}

	if err := db.WithContext(ctx).Create(&modelNotification).Error; err != nil {
		n.logger.Errorf("[NotificationRepository-1] CreateNotification: %v", err.Error())
		return err
	}

	return nil
}

// GetNotificationById implements [NotificationRepositoryInterface].
func (n *notificationRepository) GetNotificationById(ctx context.Context, notificationId uint, db *gorm.DB) (*entity.NotificationEntity, error) {
	modelNotification := model.Notification{}

	sqlMain := db.WithContext(ctx).
		Where("id = ?", notificationId).
		Omit("created_at", "updated_at", "deleted_at")

	if err := sqlMain.First(&modelNotification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New(utils.DATA_NOT_FOUND)
		}
		n.logger.Errorf("[NotificationRepository-1] GetNotificationById: %v", err.Error())
		return nil, err
	}

	return &entity.NotificationEntity{
		ID:               modelNotification.ID,
		Subject:          modelNotification.Subject,
		Status:           modelNotification.Status,
		SentAt:           modelNotification.SentAt,
		ReadAt:           modelNotification.ReadAt,
		Message:          modelNotification.Message,
		NotificationType: modelNotification.NotificationType,
	}, nil
}

// GetAllNotifications implements [NotificationRepositoryInterface].
func (n *notificationRepository) GetAllNotifications(ctx context.Context, query entity.NotificationQueryString, db *gorm.DB) ([]entity.NotificationEntity, int64, int64, error) {
	modelNotifications := []model.Notification{}

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
		n.logger.Errorf("[NotificationRepository-1] GetAllNotifications: %v", err.Error())
		return nil, 0, 0, err
	}

	orderSort := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)

	totalPages := int(math.Ceil(float64(countData) / float64(query.Limit)))
	if err := sqlMain.Order(orderSort).
		Limit(int(query.Limit)).
		Offset(int(offset)).
		Find(&modelNotifications).Error; err != nil {
		n.logger.Errorf("[NotificationRepository-2] GetAllNotifications: %v", err.Error())
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

func NewNotificationRepository(logger *log.Logger) NotificationRepositoryInterface {
	return &notificationRepository{logger: logger}
}
