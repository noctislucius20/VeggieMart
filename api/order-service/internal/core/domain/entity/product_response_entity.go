package entity

import "time"

type ProductResponseEntity struct {
	ID           int64   `json:"id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	SalePrice    float64 `json:"sale_price"`
	RegularPrice float64 `json:"regular_price"`
	Unit         string  `json:"unit"`
	Weight       int     `json:"weight"`
}

type ProductConsumerResponse struct {
	ID           int64                     `json:"id"`
	Name         string                    `json:"name"`
	Image        string                    `json:"image"`
	RegularPrice float64                   `json:"regular_price"`
	SalePrice    float64                   `json:"sale_price"`
	Unit         string                    `json:"unit"`
	Weight       int64                     `json:"weight"`
	Stock        int64                     `json:"stock"`
	CreatedAt    time.Time                 `json:"created_at"`
	Childs       []ProductConsumerResponse `json:"childs"`
}

type ProductHttpClientResponse struct {
	Message string                  `json:"message"`
	Data    []ProductResponseEntity `json:"data"`
}
