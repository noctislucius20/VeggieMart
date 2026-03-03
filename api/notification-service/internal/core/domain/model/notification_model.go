package model

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID               uint            `gorm:"primaryKey"`
	NotificationType string          `gorm:"type:varchar(50);not null"`
	ReceiverID       *int64          `gorm:"type:bigint;null"`
	ReceiverEmail    *string         `gorm:"type:varchar(80);null"`
	Subject          *string         `gorm:"type:varchar(255);null"`
	Message          string          `gorm:"type:text;not null"`
	Status           string          `gorm:"type:varchar(50);not null"`
	SentAt           *time.Time      `gorm:"type:timestamp;null"`
	ReadAt           *time.Time      `gorm:"type:timestamp;null"`
	CreatedAt        time.Time       `gorm:"autoCreateTime"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime"`
	DeletedAt        *gorm.DeletedAt `gorm:"index"`
}
