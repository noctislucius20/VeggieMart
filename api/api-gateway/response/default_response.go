package response

type DefaultResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ResponseSuccess(data any) DefaultResponse {
	return DefaultResponse{
		Message: "success",
		Data:    data,
	}
}

func ResponseFailed(message string) DefaultResponse {
	return DefaultResponse{
		Message: message,
		Data:    nil,
	}
}
