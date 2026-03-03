package model

import "time"

type OutboxEvent struct {
	ID          int64 `gorm:"primaryKey"`
	EventType   string
	AggregateID string
	Payload     string `gorm:"type:json"`
	Status      string `gorm:"index"`
	RetryCount  int64
	CreatedAt   time.Time
	NextRetryAt *time.Time
}
