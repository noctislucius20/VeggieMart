package request

type CreateOrderRequest struct {
	OrderDate    string                `json:"order_date" validate:"required,datetime=2006-01-02"`
	ShippingType string                `json:"shipping_type" validate:"required,oneof=PICKUP DELIVERY,max=20"`
	PaymentType  string                `json:"payment_type" validate:"omitempty,max=50"`
	Remarks      string                `json:"remarks" validate:"omitempty"`
	OrderTime    string                `json:"order_time" validate:"required,datetime=15:04:05"`
	OrderDetails []OrderDetailsRequest `json:"order_details" validate:"required,min=1,dive"`
}

type OrderDetailsRequest struct {
	ProductID int64 `json:"product_id" validate:"required,gte=1"`
	Quantity  int64 `json:"quantity" validate:"required,gte=1"`
}

type OrderUpdateStatusRequest struct {
	Status  string `json:"status" validate:"required,oneof=PENDING CONFIRMED PROCESSING SHIPPED DELIVERED CANCELLED,max=20"`
	Remarks string `json:"remarks" validate:"omitempty"`
}

type OrderBatchRequest struct {
	IDOrders []int64 `json:"id_orders" validate:"required,min=1"`
}
