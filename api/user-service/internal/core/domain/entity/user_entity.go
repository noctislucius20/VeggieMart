package entity

type UserEntity struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	RoleName   string `json:"role_name"`
	RoleID     int64  `json:"role_id"`
	Address    string `json:"address"`
	Lat        string `json:"lat"`
	Lng        string `json:"lng"`
	Phone      string `json:"phone"`
	Photo      string `json:"photo"`
	IsVerified bool   `json:"is_verified"`
	Token      string `json:"token"`
}

type QueryStringEntity struct {
	Search    string
	Page      int64
	Limit     int64
	OrderBy   string
	OrderType string
}

type SessionEntity struct {
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	LoggedIn  bool   `json:"logged_in"`
	CreatedAt string `json:"created_at"`
	Token     string `json:"token"`
	RoleName  string `json:"role_name"`
}
