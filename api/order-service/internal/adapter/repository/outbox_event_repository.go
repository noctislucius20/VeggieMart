package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"order-service/internal/core/domain/entity"
	"order-service/internal/core/domain/model"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OutboxEventInterface interface {
	CreateEvent(ctx context.Context, publishName string, payload any, orderId *int64, db *gorm.DB) error
	GetAllPendingEvent(ctx context.Context, db *gorm.DB) ([]entity.OutboxEventEntity, error)
	UpdateFailedEvent(ctx context.Context, outboxIds []int64, db *gorm.DB) error
	UpdatePublishedEvent(ctx context.Context, outboxIds []int64, db *gorm.DB) error
}

type outboxEventRepository struct {
	logger *log.Logger
}

// UpdateFailedEvent implements OutboxEventInterface.
func (o *outboxEventRepository) UpdateFailedEvent(ctx context.Context, outboxIds []int64, db *gorm.DB) error {
	outboxModels := []model.OutboxEvent{}

	if err := db.WithContext(ctx).Where("id IN ?", outboxIds).Find(&outboxModels).Error; err != nil {
		o.logger.Errorf("[OutboxEventRepository-1] UpdateFailedEvent: %v", err.Error())
		return err
	}

	updatedOutbox := []model.OutboxEvent{}

	for _, outbox := range outboxModels {
		if outbox.RetryCount > 10 {
			outbox.Status = "FAILED"
			outbox.NextRetryAt = nil
			outbox.RetryCount = 0
			continue
		} else {
			outbox.Status = "PENDING"
		}

		outbox.RetryCount += 1

		next := time.Now().Add(time.Duration(math.Pow(float64(outbox.RetryCount), 2)) * time.Second)
		outbox.NextRetryAt = &next

		updatedOutbox = append(updatedOutbox, outbox)
	}

	if err := db.WithContext(ctx).Save(&updatedOutbox).Error; err != nil {
		o.logger.Errorf("[OutboxEventRepository-2] UpdateFailedEvent: %v", err.Error())
		return err
	}

	return nil
}

// UpdatePublishedEvent implements OutboxEventInterface.
func (o *outboxEventRepository) UpdatePublishedEvent(ctx context.Context, outboxIds []int64, db *gorm.DB) error {
	outboxModels := []model.OutboxEvent{}

	if err := db.WithContext(ctx).Where("id IN ?", outboxIds).Find(&outboxModels).Error; err != nil {
		o.logger.Errorf("[OutboxEventRepository-1] UpdatePublishedEvent: %v", err.Error())
		return err
	}

	updatedOutbox := []model.OutboxEvent{}

	for _, outbox := range outboxModels {
		outbox.Status = "PUBLISHED"

		updatedOutbox = append(updatedOutbox, outbox)
	}

	if err := db.WithContext(ctx).Save(&updatedOutbox).Error; err != nil {
		o.logger.Errorf("[OutboxEventRepository-2] UpdatePublishedEvent: %v", err.Error())
		return err
	}

	return nil
}

// GetAllPendingEvent implements OutboxEventInterface.
func (o *outboxEventRepository) GetAllPendingEvent(ctx context.Context, db *gorm.DB) ([]entity.OutboxEventEntity, error) {
	outboxModels := []model.OutboxEvent{}
	outboxEntities := []entity.OutboxEventEntity{}

	if err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where("status = ? AND next_retry_at <= ?", "PENDING", time.Now()).
		Limit(10).
		Find(&outboxModels).Error; err != nil {
		o.logger.Errorf("[OutboxEventRepository-1] GetAllPendingEvent: %v", err.Error())
		return nil, err
	}

	if len(outboxModels) == 0 {
		return nil, nil
	}

	updatedOutbox := []model.OutboxEvent{}

	for _, val := range outboxModels {
		val.Status = "PROCESSING"

		updatedOutbox = append(updatedOutbox, val)

		outboxEntities = append(outboxEntities, entity.OutboxEventEntity{
			ID:          val.ID,
			EventType:   val.EventType,
			AggregateID: val.AggregateID,
			Payload:     val.Payload,
			Status:      val.Status,
			RetryCount:  val.RetryCount,
			NextRetryAt: val.NextRetryAt,
		})
	}

	if err := db.WithContext(ctx).Save(&updatedOutbox).Error; err != nil {
		o.logger.Errorf("[OutboxEventRepository-2] GetAllPendingEvent: %v", err.Error())
		return nil, err
	}

	return outboxEntities, nil

}

// CreateEvent implements OutboxEventInterface.
func (o *outboxEventRepository) CreateEvent(ctx context.Context, publishName string, payload any, orderId *int64, db *gorm.DB) error {
	parsedPayload, _ := json.Marshal(payload)

	timeNow := time.Now()

	aggregateId := string("")
	if orderId != nil {
		aggregateId = fmt.Sprintf("%d", *orderId)
	} else {
		aggregateId = ""
	}

	outboxModel := model.OutboxEvent{
		EventType:   publishName,
		AggregateID: aggregateId,
		Payload:     string(parsedPayload),
		Status:      "PENDING",
		NextRetryAt: &timeNow,
	}

	if err := db.WithContext(ctx).Create(&outboxModel).Error; err != nil {
		o.logger.Errorf("[OutboxEventRepository-1] CreateEvent: %v", err.Error())
		return err
	}

	return nil
}

func NewOutboxEventRepository(logger *log.Logger) OutboxEventInterface {
	return &outboxEventRepository{
		logger: logger,
	}
}
