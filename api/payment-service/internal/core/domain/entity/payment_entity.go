package entity

type PaymentEntity struct {
	ID               int64              `json:"id"`
	UserID           int64              `json:"user_id"`
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
	ID            int64  `json:"id"`
	OrderCode     string `json:"order_code"`
	OrderDatetime string `json:"order_datetime"`
	Status        string `json:"status"`
	ShippingFee   int64  `json:"shipping_fee"`
	ShippingType  string `json:"shipping_type"`
	Remarks       string `json:"remarks"`
	TotalAmount   int64  `json:"total_amount"`
}

type CustomerEntity struct {
	CustomerID      int64  `json:"customer_id"`
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	CustomerAddress string `json:"customer_address"`
	CustomerEmail   string `json:"customer_email"`
}
