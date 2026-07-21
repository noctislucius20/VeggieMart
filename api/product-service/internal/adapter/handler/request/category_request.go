package request

type CategoryRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Icon        string `json:"icon" validate:"required,url,max=255"`
	Description string `json:"description" validate:"omitempty"`
	Status      string `json:"status" validate:"required,oneof=true false"`
	ParentID    *int64 `json:"parent_id" validate:"omitempty"`
}
