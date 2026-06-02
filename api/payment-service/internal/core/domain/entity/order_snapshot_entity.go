package entity

import (
	"time"
)

type OrderSnapshotEntity struct {
	ID              int64      `json:"id"`
	OrderID         int64      `json:"order_id"`
	OrderCode       string     `json:"order_code"`
	OrderDatetime   time.Time  `json:"order_datetime"`
	Status          string     `json:"status"`
	PaymentMethod   string     `json:"payment_method"`
	ShippingFee     float64    `json:"shipping_fee"`
	ShippingType    string     `json:"shipping_type"`
	Remarks         string     `json:"remarks"`
	TotalAmount     float64    `json:"total_amount"`
	CustomerID      int64      `json:"customer_id"`
	CustomerName    string     `json:"customer_name"`
	CustomerEmail   string     `json:"customer_email"`
	CustomerAddress string     `json:"customer_address"`
	CustomerPhone   string     `json:"customer_phone"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastUsed        *time.Time `json:"last_used"`
}
