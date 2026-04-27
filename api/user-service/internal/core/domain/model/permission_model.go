package model

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID        int64           `gorm:"primaryKey"`
	Resource  string          `gorm:"type:varchar(255);not null;uniqueIndex:idx_permission_unique,where:deleted_at IS NULL"`
	Action    string          `gorm:"type:varchar(255);not null;uniqueIndex:idx_permission_unique,where:deleted_at IS NULL"`
	Scope     string          `gorm:"type:varchar(255);not null;uniqueIndex:idx_permission_unique,where:deleted_at IS NULL"`
	CreatedAt time.Time       `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt time.Time       `gorm:"type:timestamp"`
	DeletedAt *gorm.DeletedAt `gorm:"index"`
	Roles     []Role          `gorm:"many2many:role_permission"`
}
