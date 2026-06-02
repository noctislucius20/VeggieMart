package entity

type OrderHttpClientResponseList struct {
	Message string                      `json:"message"`
	Data    []OrderDetailResponseEntity `json:"data"`
}

type OrderDetailResponseEntity struct {
	ID            int64         `json:"id"`
	OrderCode     string        `json:"order_code"`
	OrderDatetime string        `json:"order_datetime"`
	Status        string        `json:"status"`
	TotalAmount   int64         `json:"total_amount"`
	ShippingFee   int64         `json:"shipping_fee"`
	ShippingType  string        `json:"shipping_type"`
	Remarks       string        `json:"remarks"`
	Customer      OrderCustomer `json:"customer"`
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
		OrderID int64 `json:"order_id"`
	} `json:"data"`
}

type OrderConsumerResponse struct {
	ID            int64         `json:"id"`
	OrderCode     string        `json:"order_code"`
	OrderDatetime string        `json:"order_datetime"`
	Status        string        `json:"status"`
	PaymentMethod string        `json:"payment_method"`
	ShippingFee   int64         `json:"shipping_fee"`
	ShippingType  string        `json:"shipping_type"`
	Remarks       string        `json:"remarks"`
	TotalAmount   int64         `json:"total_amount"`
	Customer      OrderCustomer `json:"customer"`
}
