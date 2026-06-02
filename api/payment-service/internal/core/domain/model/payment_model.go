package model

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID               int64           `gorm:"primaryKey" json:"id"`
	UserID           int64           `gorm:"type:bigint;not null;index" json:"user_id"`
	PaymentMethod    string          `gorm:"type:varchar(50);not null" json:"payment_method"`
	PaymentStatus    string          `gorm:"type:varchar(50);not null" json:"payment_status"`
	PaymentGatewayID *string         `gorm:"type:varchar(50)" json:"payment_gateway_id,omitempty"`
	GrossAmount      float64         `gorm:"type:decimal(10,2);not null" json:"gross_amount"`
	PaymentURL       *string         `gorm:"type:text" json:"payment_url,omitempty"`
	CreatedAt        time.Time       `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"type:timestamp" json:"updated_at"`
	DeletedAt        *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	PaymentLogs      []PaymentLog    `gorm:"foreignKey:PaymentID;constraint:OnDelete:CASCADE"`
	OrderSnapshot    *OrderSnapshot  `gorm:"foreignKey:PaymentID;constraint:OnDelete:CASCADE"`
}
