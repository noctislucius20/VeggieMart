package model

import (
	"time"
)

type OrderSnapshot struct {
	ID              int64      `gorm:"primaryKey"`
	PaymentID       int64      `gorm:"type:bigint;not null;index"`
	OrderID         int64      `gorm:"type:bigint;not null;index"`
	OrderCode       string     `gorm:"type:varchar(64);not null"`
	OrderDatetime   time.Time  `gorm:"type:timestamp;not null"`
	Status          string     `gorm:"type:varchar(20);not null"`
	PaymentMethod   string     `gorm:"type:varchar(50);not null"`
	ShippingFee     float64    `gorm:"type:decimal(10,2);not null;default:0"`
	ShippingType    string     `gorm:"type:varchar(20);not null"`
	Remarks         string     `gorm:"type:text"`
	TotalAmount     float64    `gorm:"type:decimal(10,2);not null;default:0"`
	CustomerID      int64      `gorm:"type:bigint;not null;index"`
	CustomerName    string     `gorm:"type:varchar(255);not null"`
	CustomerEmail   string     `gorm:"type:varchar(255);not null"`
	CustomerAddress string     `gorm:"type:text"`
	CustomerPhone   string     `gorm:"type:varchar(17)"`
	CreatedAt       time.Time  `gorm:"type:timestamp"`
	UpdatedAt       time.Time  `gorm:"type:timestamp"`
	LastUsed        *time.Time `gorm:"type:timestamp"`
	Payment         Payment    `gorm:"foreignKey:PaymentID;constraint:OnDelete:CASCADE"`
}
