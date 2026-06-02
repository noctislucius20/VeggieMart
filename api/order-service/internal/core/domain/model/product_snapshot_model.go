package model

import (
	"time"
)

type ProductSnapshot struct {
	ID           int64     `gorm:"primaryKey"`
	OrderID      int64     `gorm:"type:bigint;not null;index"`
	ProductID    int64     `gorm:"type:bigint;not null;index"`
	Name         string    `gorm:"type:varchar(100);not null"`
	Image        string    `gorm:"type:varchar(255);not null"`
	RegularPrice float64   `gorm:"type:bigint;default:0"`
	SalePrice    float64   `gorm:"type:bigint;default:0"`
	Unit         string    `gorm:"type:varchar(120);default:gram"`
	Weight       int64     `gorm:"type:bigint;default:0"`
	CreatedAt    time.Time `gorm:"type:timestamp"`
	LastUsed     time.Time `gorm:"type:timestamp"`
	Order        Order     `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}
