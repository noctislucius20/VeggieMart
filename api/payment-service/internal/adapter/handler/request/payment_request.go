package request

type PaymentRequest struct {
	OrderID       uint    `json:"order_id" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	GrassAmount   float64 `json:"gross_amount" validate:"required"`
	UserID        uint    `json:"user_id" validate:"required"`
	Remarks       string  `json:"remarks"`
}
