package entity

type UserResponseEntity struct {
	// RoleName string `json:"role,omitempty"`
	// RoleId   int64  `json:"role_id"`
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	// Photo    string `json:"photo"`
}

type UsersHttpClientResponse struct {
	Message string               `json:"message"`
	Data    []UserResponseEntity `json:"data"`
}

type UserHttpClientResponse struct {
	Message string             `json:"message"`
	Data    UserResponseEntity `json:"data"`
}
