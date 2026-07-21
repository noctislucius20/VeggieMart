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
	BuyerID       int64  `json:"buyer_id"`
	OrderDate     string `json:"order_date"`
	OrderDatetime string `json:"order_datetime"`
	PaymentMethod string `json:"payment_method"`
	Status        string `json:"status"`
	TotalAmount   int64  `json:"total_amount"`
	ShippingType  string `json:"shipping_type"`
	ShippingFee   int64  `json:"shipping_fee"`
	OrderTime     string `json:"order_time"`
	Remarks       string `json:"remarks"`
	BuyerName     string `json:"buyer_name"`
	BuyerEmail    string `json:"buyer_email"`
	BuyerPhone    string `json:"buyer_phone"`
	BuyerAddress  string `json:"buyer_address"`
	BuyerLat      string `json:"buyer_lat"`
	BuyerLng      string `json:"buyer_lng"`
}

type CustomerEntity struct {
	CustomerID      int64  `json:"customer_id"`
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	CustomerAddress string `json:"customer_address"`
	CustomerEmail   string `json:"customer_email"`
}
