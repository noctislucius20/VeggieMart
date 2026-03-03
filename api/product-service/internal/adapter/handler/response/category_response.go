package response

type CategoryListResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	Slug         string `json:"slug"`
	Status       string `json:"status"`
	TotalProduct int    `json:"total_product"`
}

type CategoryResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Slug        string `json:"slug"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type CategoryHomeListResponse struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Slug string `json:"slug"`
}

type CategoryShopListResponse struct {
	Name   string                     `json:"name"`
	Slug   string                     `json:"slug"`
	Childs []CategoryShopListResponse `json:"childs"`
}
