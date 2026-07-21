package request

type CartRequest struct {
	ProductID int64 `json:"product_id" validate:"required,gte=1"`
	Quantity  int64 `json:"quantity" validate:"required,gte=1"`
}
