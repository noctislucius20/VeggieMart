package entity

import "time"

type ProductEntity struct {
	ID           int64           `json:"id"`
	CategoryID   int64           `json:"category_id,omitempty"`
	CategorySlug string          `json:"category_slug"`
	ParentID     *int64          `json:"parent_id"`
	Name         string          `json:"name"`
	Image        string          `json:"image"`
	Description  string          `json:"description"`
	RegularPrice float64         `json:"regular_price"`
	SalePrice    float64         `json:"sale_price"`
	Unit         string          `json:"unit"`
	Weight       int64           `json:"weight"`
	Stock        int64           `json:"stock"`
	Variant      int64           `json:"variant"`
	Status       string          `json:"status"`
	CategoryName string          `json:"category_name"`
	CreatedAt    time.Time       `json:"created_at"`
	Childs       []ProductEntity `json:"childs"`
}

type QueryStringProduct struct {
	Search       string
	Page         int64
	Limit        int64
	OrderBy      string
	OrderType    string
	CategorySlug string
	StartPrice   int64
	EndPrice     int64
}
