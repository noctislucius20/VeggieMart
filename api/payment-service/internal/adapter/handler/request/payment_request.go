package request

type PaymentRequest struct {
	OrderID       int64   `json:"order_id" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	GrassAmount   float64 `json:"gross_amount" validate:"required"`
	Remarks       string  `json:"remarks"`
}
