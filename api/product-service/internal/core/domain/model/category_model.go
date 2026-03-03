package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          int64           `gorm:"primaryKey"`
	ParentID    *int64          `gorm:"type:bigint"`
	Name        string          `gorm:"type:varchar(100);not null"`
	Icon        string          `gorm:"type:varchar(255);not null"`
	Status      bool            `gorm:"type:boolean;default:true;index:idx_categories_status"`
	Slug        string          `gorm:"type:varchar(120);not null;uniqeIndex:idx_categories_slug_unique,where:deleted_at IS NULL"`
	Description string          `gorm:"type:text"`
	CreatedAt   time.Time       `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt   time.Time       `gorm:"type:timestamp"`
	DeletedAt   *gorm.DeletedAt `gorm:"index"`
	Products    []Product       `gorm:"foreignKey:CategoryID"`
}

type CategoryDTO struct {
	CategoryID int64
	ProductID  int64
	Slug       string
}
