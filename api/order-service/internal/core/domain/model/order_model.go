package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID           int64           `gorm:"primaryKey"`
	OrderCode    string          `gorm:"type:varchar(64);not null;uniqueIndex:idx_orders_order_code_unique,where:deleted_at IS NULL"`
	BuyerID      int64           `gorm:"type:bigint;not null"`
	OrderDate    time.Time       `gorm:"type:date;default:current_timestamp"`
	Status       string          `gorm:"type:varchar(20);not null;default:PENDING;index:idx_orders_status"`
	TotalAmount  float64         `gorm:"type:decimal(10,2);not null;default:0"`
	ShippingType string          `gorm:"type:varchar(20);not null;default:PICKUP"`
	ShippingFee  float64         `gorm:"type:decimal(10,2);not null;default:0"`
	OrderTime    time.Time       `gorm:"type:time"`
	Remarks      string          `gorm:"type:text"`
	CreatedAt    time.Time       `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt    time.Time       `gorm:"type:timestamp"`
	DeletedAt    *gorm.DeletedAt `gorm:"index"`
	OrderItems   []OrderItem     `gorm:"foreignKey:OrderID"`
}
