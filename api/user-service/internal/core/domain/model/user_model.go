package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         int64           `gorm:"primaryKey"`
	Name       string          `gorm:"type:varchar(255);not null"`
	Email      string          `gorm:"type:varchar(255);not null;uniqueIndex:idx_users_email_unique,where:deleted_at IS NULL"`
	Password   string          `gorm:"type:varchar(255);not null"`
	Address    string          `gorm:"type:text"`
	Phone      string          `gorm:"type:varchar(17)"`
	Photo      string          `gorm:"type:varchar(255)"`
	Lat        string          `gorm:"type:varchar(50)"`
	Lng        string          `gorm:"type:varchar(50)"`
	IsVerified bool            `gorm:"type:boolean;default:false;index:idx_users_is_verified"`
	CreatedAt  time.Time       `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt  time.Time       `gorm:"type:timestamp"`
	DeletedAt  *gorm.DeletedAt `gorm:"index"`
	Roles      []Role          `gorm:"many2many:user_role"`
}
