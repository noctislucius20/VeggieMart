package request

type PaymentRequest struct {
	OrderID       int64  `json:"order_id" validate:"required,gte=1"`
	PaymentMethod string `json:"payment_method" validate:"required,oneof=COD TRANSFER,max=50"`
	Remarks       string `json:"remarks" validate:"omitempty"`
}
