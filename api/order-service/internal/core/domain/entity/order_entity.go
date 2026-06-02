package entity

type OrderEntity struct {
	ID               int64                   `json:"id"`
	OrderCode        string                  `json:"order_code"`
	BuyerID          int64                   `json:"buyer_id"`
	OrderDate        string                  `json:"order_date"`
	PaymentMethod    string                  `json:"payment_method"`
	Status           string                  `json:"status"`
	TotalAmount      int64                   `json:"total_amount"`
	ShippingType     string                  `json:"shipping_type"`
	ShippingFee      int64                   `json:"shipping_fee"`
	OrderTime        string                  `json:"order_time"`
	Remarks          string                  `json:"remarks"`
	BuyerName        string                  `json:"buyer_name"`
	BuyerEmail       string                  `json:"buyer_email"`
	BuyerPhone       string                  `json:"buyer_phone"`
	BuyerAddress     string                  `json:"buyer_address"`
	BuyerLat         string                  `json:"buyer_lat"`
	BuyerLng         string                  `json:"buyer_lng"`
	OrderItems       []OrderItemEntity       `json:"order_items"`
	ProductSnapshots []ProductSnapshotEntity `json:"product_snapshots"`
}

type OrderQueryString struct {
	Search  string
	Page    int64
	Limit   int64
	Status  string
	BuyerID int64
}
