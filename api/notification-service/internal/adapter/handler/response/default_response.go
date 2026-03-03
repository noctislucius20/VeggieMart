package response

type DefaultResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type DefaultResponseWithPaginations struct {
	Message    string      `json:"message"`
	Data       any         `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page       int64 `json:"page"`
	TotalCount int64 `json:"total_count"`
	PerPage    int64 `json:"per_page"`
	TotalPage  int64 `json:"total_page"`
}

func ResponseSuccess(data any) DefaultResponse {
	return DefaultResponse{
		Message: "success",
		Data:    data,
	}
}

func ResponseWithPaginationsSuccess(data any, pagination Pagination) DefaultResponseWithPaginations {
	return DefaultResponseWithPaginations{
		Message:    "success",
		Data:       data,
		Pagination: &pagination,
	}
}

func ResponseFailed(message string) DefaultResponse {
	return DefaultResponse{
		Message: message,
		Data:    nil,
	}
}
