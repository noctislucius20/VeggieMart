package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           int64           `gorm:"primaryKey"`
	CategoryID   int64           `gorm:"not null;index"`
	ParentID     *int64          `gorm:"type:bigint"`
	Name         string          `gorm:"type:varchar(100);not null"`
	Image        string          `gorm:"type:varchar(255);not null"`
	Description  string          `gorm:"type:text"`
	RegularPrice float64         `gorm:"type:bigint;default:0"`
	SalePrice    float64         `gorm:"type:bigint;default:0;index:idx_products_sale_price"`
	Unit         string          `gorm:"type:varchar(120);default:gram"`
	Weight       int64           `gorm:"type:bigint;default:0"`
	Stock        int64           `gorm:"type:bigint;default:0"`
	Variant      int64           `gorm:"type:bigint;default:1"`
	Status       string          `gorm:"type:varchar(120);default:DRAFT;index:idx_products_status"`
	CreatedAt    time.Time       `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt    time.Time       `gorm:"type:timestamp"`
	DeletedAt    *gorm.DeletedAt `gorm:"index"`
	Childs       []Product       `gorm:"foreignkey:ParentID"`
	Categories   Category        `gorm:"constraint:OnDelete:CASCADE;foreignkey:CategoryID"`
}
