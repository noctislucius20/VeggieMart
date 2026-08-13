package model

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID                 int64           `gorm:"primaryKey"`
	NotificationType   string          `gorm:"type:varchar(50);not null"`
	NotificationTypeID *int64          `gorm:"type:bigint"`
	NotificationMethod string          `gorm:"type:varchar(50);not null"`
	ReceiverID         *int64          `gorm:"type:bigint"`
	ReceiverEmail      *string         `gorm:"type:varchar(80)"`
	Subject            *string         `gorm:"type:varchar(255)"`
	Message            string          `gorm:"type:text;not null"`
	Status             string          `gorm:"type:varchar(50);not null"`
	SentAt             *time.Time      `gorm:"type:timestamp"`
	ReadAt             *time.Time      `gorm:"type:timestamp"`
	CreatedAt          time.Time       `gorm:"autoCreateTime"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime"`
	DeletedAt          *gorm.DeletedAt `gorm:"index"`
}
