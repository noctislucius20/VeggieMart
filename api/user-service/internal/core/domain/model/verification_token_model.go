package model

import (
	"time"

	"gorm.io/gorm"
)

type VerificationToken struct {
	ID        int64          `gorm:"primaryKey"`
	UserID    int64          `gorm:"not null;index"`
	Token     string         `gorm:"type:varchar(255);not null;index:idx_verification_tokens_token"`
	TokenType string         `gorm:"type:varchar(20);not null"`
	ExpiresAt time.Time      `gorm:"type:timestamp;not null"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt *time.Time     `gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      User           `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
}

type VerificationTokenDTO struct {
	ID        int64
	UserID    int64
	TokenType string
	ExpiresAt time.Time
	UserEmail string
	UserName  string
	RoleName  string
}
