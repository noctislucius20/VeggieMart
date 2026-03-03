package entity

import "time"

type OutboxEventEntity struct {
	ID          int64
	EventType   string
	AggregateID string
	Payload     string
	Status      string
	RetryCount  int64
	NextRetryAt *time.Time
}
