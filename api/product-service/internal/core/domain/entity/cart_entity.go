package entity

type CartItem struct {
	ProductID     int64         `json:"product_id"`
	Quantity      int64         `json:"quantity"`
	ProductDetail ProductEntity `json:"product_detail"`
}
