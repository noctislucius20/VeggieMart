package request

type ProductRequest struct {
	ProductName        string                 `json:"product_name" validate:"required,min=3,max=100"`
	CategorySlug       string                 `json:"category_slug" validate:"required,max=120"`
	Unit               string                 `json:"unit" validate:"required,max=120"`
	Variant            int64                  `json:"variant" validate:"required,gte=1"`
	ProductDescription string                 `json:"product_description" validate:"required"`
	Status             string                 `json:"status" validate:"required,oneof=DRAFT ACTIVE INACTIVE,max=120"`
	VariantDetail      []ProductDetailRequest `json:"variant_detail" validate:"required,min=1,dive"`
}

type ProductDetailRequest struct {
	Stock        int64  `json:"stock" validate:"required,gte=0"`
	ProductImage string `json:"product_image" validate:"required,url,max=255"`
	Weight       int64  `json:"weight" validate:"required,gte=0"`
	SalePrice    int64  `json:"sale_price" validate:"required,gte=0"`
	RegularPrice int64  `json:"regular_price" validate:"required,gte=0"`
}

type ProductUpdateStockRequest struct {
	ProductID int64 `json:"product_id" validate:"required,gte=1"`
	Quantity  int64 `json:"quantity" validate:"required"`
}

type ProductBatchRequest struct {
	IDProducts []int64 `json:"id_products" validate:"required,min=1"`
}
