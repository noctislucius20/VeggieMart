package response

import "time"

type ProductListResponse struct {
	ID            int64     `json:"id"`
	ProductName   string    `json:"product_name"`
	ParentID      *int64    `json:"parent_id"`
	ProductImage  string    `json:"product_image"`
	CategoryName  string    `json:"category_name"`
	ProductStatus string    `json:"product_status"`
	SalePrice     int64     `json:"sale_price"`
	CreatedAt     time.Time `json:"created_at"`
}

type ProductDetailResponse struct {
	ID            int64                  `json:"id"`
	ProductName   string                 `json:"product_name"`
	ParentID      *int64                 `json:"parent_id"`
	ProductImage  string                 `json:"product_image"`
	CategoryName  string                 `json:"category_name"`
	ProductStatus string                 `json:"product_status"`
	SalePrice     int64                  `json:"sale_price"`
	RegularPrice  int64                  `json:"regular_price"`
	CreatedAt     time.Time              `json:"created_at"`
	Unit          string                 `json:"unit"`
	Weight        int64                  `json:"weight"`
	Stock         int64                  `json:"stock"`
	Childs        []ProductChildResponse `json:"child"`
}

type ProductChildResponse struct {
	ID           int64 `json:"id"`
	Weight       int64 `json:"weight"`
	Stock        int64 `json:"stock"`
	RegularPrice int64 `json:"regular_price"`
	SalePrice    int64 `json:"sale_price"`
}

type ProductHomeListResponse struct {
	ID           int64  `json:"id"`
	ProductName  string `json:"product_name"`
	ProductImage string `json:"product_image"`
	CategoryName string `json:"category_name"`
	SalePrice    int64  `json:"sale_price"`
	RegularPrice int64  `json:"regular_price"`
}

type ProductHomeDetailResponse struct {
	ID           int64                      `json:"id"`
	ProductName  string                     `json:"product_name"`
	CategoryName string                     `json:"category_name"`
	Description  string                     `json:"description"`
	Unit         string                     `json:"unit"`
	Childs       []ProductHomeChildResponse `json:"childs"`
}

type ProductHomeChildResponse struct {
	ID           int64  `json:"id"`
	Weight       int64  `json:"weight"`
	Stock        int64  `json:"stock"`
	RegularPrice int64  `json:"regular_price"`
	SalePrice    int64  `json:"sale_price"`
	Image        string `json:"image"`
}

type ProductBatchResponse struct {
	ID           int64  `json:"id"`
	ProductImage string `json:"product_image"`
	ProductName  string `json:"product_name"`
	SalePrice    int64  `json:"sale_price"`
	Weight       int64  `json:"weight"`
	Unit         string `json:"unit"`
}
