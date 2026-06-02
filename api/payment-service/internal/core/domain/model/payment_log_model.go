package model

import "time"

type PaymentLog struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	PaymentID int64     `gorm:"not null;index" json:"payment_id"`
	Status    string    `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
}
