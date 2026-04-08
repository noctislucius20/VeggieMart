package model

import "time"

type OutboxEvent struct {
	ID          int64      `gorm:"primaryKey"`
	EventType   string     `gorm:"type:varchar(100)"`
	AggregateID string     `gorm:"type:varchar(100)"`
	Payload     string     `gorm:"type:json"`
	Status      string     `gorm:"type:string;index:idx_outbox_events_status"`
	RetryCount  int64      `gorm:"type:int"`
	CreatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp"`
	NextRetryAt *time.Time `gorm:"type:timestamp"`
}
