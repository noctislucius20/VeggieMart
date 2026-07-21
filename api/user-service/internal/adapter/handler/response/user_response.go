package response

type SignInResponse struct {
	AccessToken string `json:"access_token"`
	ID          int64  `json:"id"`
	Role        string `json:"role"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Photo       string `json:"photo"`
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

type UpdateProfileResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Role        string `json:"role"`
	Lat         string `json:"lat"`
	Lng         string `json:"lng"`
	Photo       string `json:"photo"`
	AccessToken string `json:"access_token"`
}

type CustomerResponse struct {
	RoleName string `json:"role,omitempty"`
	RoleID   int64  `json:"role_id"`
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
