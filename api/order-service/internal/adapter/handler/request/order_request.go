package request

type CreateOrderRequest struct {
	BuyerID      int64                 `json:"buyer_id" validate:"required"`
	OrderDate    string                `json:"order_date" validate:"required"`
	TotalAmount  int64                 `json:"total_amount" validate:"required"`
	ShippingType string                `json:"shipping_type" validate:"required"`
	PaymentType  string                `json:"payment_type"`
	Remarks      string                `json:"remarks"`
	OrderTime    string                `json:"order_time" validate:"required"`
	OrderDetails []OrderDetailsRequest `json:"order_details" validate:"required"`
}

type OrderDetailsRequest struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Quantity  int64 `json:"quantity" validate:"required"`
}

type OrderUpdateStatusRequest struct {
	Status  string `json:"status" validate:"required"`
	Remarks string `json:"remarks"`
}

type OrderBatchRequest struct {
	IDOrders []int64 `json:"id_orders" validate:"required"`
}
