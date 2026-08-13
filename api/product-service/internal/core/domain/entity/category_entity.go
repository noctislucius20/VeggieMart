package entity

type CategoryEntity struct {
	ID          int64           `json:"id"`
	ParentID    *int64          `json:"parent_id"`
	Name        string          `json:"name"`
	Icon        string          `json:"icon"`
	Status      string          `json:"status"`
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	Products    []ProductEntity `json:"products"`
}

type QueryStringEntity struct {
	Search    string
	Page      int64
	Limit     int64
	OrderBy   string
	OrderType string
	Status    string
}
