package entity

type PaymentEntity struct {
	ID               uint               `json:"id"`
	OrderID          uint               `json:"order_id"`
	UserID           uint               `json:"user_id"`
	PaymentMethod    string             `json:"payment_method"`
	PaymentStatus    string             `json:"payment_status"`
	PaymentGatewayID string             `json:"payment_gateway_id"`
	PaymentAt        string             `json:"payment_at"`
	GrossAmount      float64            `json:"gross_amount"`
	PaymentURL       string             `json:"payment_url"`
	Remarks          string             `json:"remarks"`
	Order            OrderEntity        `json:"order"`
	Customer         CustomerEntity     `json:"customer"`
	PaymentLogs      []PaymentLogEntity `json:"payment_logs"`
}

type QueryStringPayment struct {
	Limit     int64
	Page      int64
	UserID    int64
	Status    string
	OrderType string
	OrderBy   string
	Search    string
}

type OrderEntity struct {
	OrderCode         string `json:"order_code"`
	OrderShippingType string `json:"order_shipping_type"`
	OrderAt           string `json:"order_at"`
	OrderStatus       string `json:"order_status"`
	OrderRemarks      string `json:"order_remarks"`
}

type CustomerEntity struct {
	CustomerName    string `json:"customer_name"`
	CustomerEmail   string `json:"customer_email"`
	CustomerAddress string `json:"customer_address"`
}
