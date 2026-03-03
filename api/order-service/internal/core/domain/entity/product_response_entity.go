package entity

type ProductResponseEntity struct {
	ID          int64  `json:"id"`
	ProductName string `json:"product_name"`
	// ParentID      int                          `json:"parent_id"`
	ProductImage string `json:"product_image"`
	// CategoryName  string                       `json:"category_name"`
	// ProductStatus string                       `json:"product_status"`
	SalePrice float64 `json:"sale_price"`
	// RegulerPrice  float64                      `json:"reguler_price"`
	// CreatedAt     time.Time                    `json:"created_at"`
	Unit   string `json:"unit"`
	Weight int    `json:"weight"`
	// Stock         int                          `json:"stock"`
	// Child         []ChildProductResponseEntity `json:"child"`
}

type ProductHttpClientResponse struct {
	Message string                  `json:"message"`
	Data    []ProductResponseEntity `json:"data"`
}
