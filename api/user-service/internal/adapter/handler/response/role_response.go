package response

type RoleResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UsersResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
