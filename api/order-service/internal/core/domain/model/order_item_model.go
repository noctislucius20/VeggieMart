package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
	ID        int64           `gorm:"primaryKey"`
	OrderID   int64           `gorm:"type:bigint;not null;index"`
	ProductID int64           `gorm:"type:bigint;not null;index"`
	Quantity  int64           `gorm:"type:int;not null;default:1"`
	CreatedAt time.Time       `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt time.Time       `gorm:"type:timestamp"`
	DeletedAt *gorm.DeletedAt `gorm:"index"`
	Order     Order           `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}
