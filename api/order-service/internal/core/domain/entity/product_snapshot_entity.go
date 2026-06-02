package entity

import "time"

type ProductSnapshotEntity struct {
	ID           int64                   `json:"id"`
	ProductID    int64                   `json:"product_id"`
	Name         string                  `json:"name"`
	Image        string                  `json:"image"`
	RegularPrice float64                 `json:"regular_price"`
	SalePrice    float64                 `json:"sale_price"`
	Unit         string                  `json:"unit"`
	Weight       int64                   `json:"weight"`
	Stock        int64                   `json:"stock"`
	CreatedAt    time.Time               `json:"created_at"`
	LastUsed     time.Time               `json:"last_used"`
	Childs       []ProductSnapshotEntity `json:"childs"`
}
