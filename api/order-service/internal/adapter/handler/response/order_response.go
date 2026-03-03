package response

type OrderListResponse struct {
	ID            int64  `json:"id"`
	OrderCode     string `json:"order_code"`
	ProductImage  string `json:"product_image"`
	CustomerName  string `json:"customer_name"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	TotalAmount   int64  `json:"total_amount"`
}

type OrderDetailResponse struct {
	ID            int64              `json:"id"`
	OrderCode     string             `json:"order_code"`
	ProductImage  string             `json:"product_image"`
	OrderDatetime string             `json:"order_datetime"`
	Status        string             `json:"status"`
	PaymentMethod string             `json:"payment_method"`
	TotalAmount   int64              `json:"total_amount"`
	ShippingFee   int64              `json:"shipping_fee"`
	ShippingType  string             `json:"shipping_type"`
	Remarks       string             `json:"remarks"`
	Customer      OrderCustomer      `json:"customer"`
	OrderItems    []OrderItemsDetail `json:"order_items"`
}

type OrderCustomer struct {
	CustomerID      int64  `json:"customer_id"`
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	CustomerAddress string `json:"customer_address"`
	CustomerEmail   string `json:"customer_email"`
}

type OrderItemsDetail struct {
	ProductName  string `json:"product_name"`
	ProductImage string `json:"product_image"`
	ProductPrice int64  `json:"product_price"`
	Quantity     int64  `json:"quantity"`
}

type OrderCustomerList struct {
	ID            int64  `json:"id"`
	OrderCode     string `json:"order_code"`
	ProductImage  string `json:"product_image"`
	ProductName   string `json:"product_name"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	TotalAmount   int64  `json:"total_amount"`
	Weight        int64  `json:"weight"`
	Unit          string `json:"unit"`
	Quantity      int64  `json:"quantity"`
	OrderDateTime string `json:"order_datetime"`
}

type OrderBatchResponse struct {
	ID           int64         `json:"id"`
	OrderCode    string        `json:"order_code"`
	ShippingType string        `json:"shipping_type"`
	Customer     OrderCustomer `json:"customer"`
}
