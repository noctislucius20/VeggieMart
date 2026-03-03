package response

type SignInResponse struct {
	AccessToken string `json:"access_token"`
	Role        string `json:"role"`
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Lat         string `json:"lat"`
	Lng         string `json:"lng"`
}

type ProfileResponse struct {
	RoleName string `json:"role"`
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Lat      string `json:"lat"`
	Lng      string `json:"lng"`
	Address  string `json:"address"`
	Photo    string `json:"photo"`
}

type CustomerResponse struct {
	RoleName string `json:"role,omitempty"`
	RoleId   int64  `json:"role_id"`
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Lat      string `json:"lat"`
	Lng      string `json:"lng"`
	Address  string `json:"address"`
	Photo    string `json:"photo"`
}

type CustomerResponseList struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Photo string `json:"photo"`
}

type CustomerBatchResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}
