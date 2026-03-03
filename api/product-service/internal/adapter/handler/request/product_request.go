package request

type ProductRequest struct {
	ProductName        string                 `json:"product_name" validate:"required"`
	CategorySlug       string                 `json:"category_slug" validate:"required"`
	Unit               string                 `json:"unit" validate:"required"`
	Variant            int64                  `json:"variant" validate:"required"`
	ProductDescription string                 `json:"product_description" validate:"required"`
	Status             string                 `json:"status" validate:"required"`
	VariantDetail      []ProductDetailRequest `json:"variant_detail" validate:"required"`
}

type ProductDetailRequest struct {
	Stock        int64  `json:"stock" validate:"required,number"`
	ProductImage string `json:"product_image" validate:"required,url"`
	Weight       int64  `json:"weight" validate:"required,number"`
	SalePrice    int64  `json:"sale_price" validate:"required,number"`
	RegularPrice int64  `json:"regular_price" validate:"required,number"`
}

type ProductBatchRequest struct {
	IDProducts []int64 `json:"id_products" validate:"required"`
}
