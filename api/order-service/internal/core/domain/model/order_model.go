package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID           int64     `gorm:"primaryKey"`
	OrderCode    string    `gorm:"order_code"`
	BuyerID      int64     `gorm:"buyer_id"`
	OrderDate    time.Time `gorm:"order_date;type:date"`
	Status       string    `gorm:"status"`
	TotalAmount  float64   `gorm:"total_amount"`
	ShippingType string    `gorm:"shipping_type"`
	ShippingFee  float64   `gorm:"shipping_fee"`
	OrderTime    string    `gorm:"order_time;type:time"`
	Remarks      string    `gorm:"remarks"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *gorm.DeletedAt
	OrderItems   []OrderItem `gorm:"foreignKey:OrderID;references:ID"`
}
