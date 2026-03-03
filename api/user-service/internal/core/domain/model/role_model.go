package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        int64           `gorm:"primaryKey"`
	Name      string          `gorm:"type:varchar(255);not null;uniqueIndex:idx_roles_name_unique,where:deleted_at IS NULL"`
	CreatedAt time.Time       `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt time.Time       `gorm:"type:timestamp"`
	DeletedAt *gorm.DeletedAt `gorm:"index"`
	Users     []User          `gorm:"many2many:user_role"`
}

type RoleDeleteDTO struct {
	RoleID         int64
	UserRoleRoleID int64
}
