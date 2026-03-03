package entity

type OrderHttpClientResponseList struct {
	Message string                      `json:"message"`
	Data    []OrderDetailResponseEntity `json:"data"`
}

type OrderDetailResponseEntity struct {
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

type OrderHttpClientResponse struct {
	Message string                    `json:"message"`
	Data    OrderDetailResponseEntity `json:"data"`
}

type OrderIDHttpResponseEntity struct {
	Message string `json:"message"`
	Data    struct {
		OrderID uint `json:"order_id"`
	} `json:"data"`
}
