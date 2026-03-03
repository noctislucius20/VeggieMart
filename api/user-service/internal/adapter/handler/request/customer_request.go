package request

type CustomerRequest struct {
	Name                 string  `json:"name" validate:"required,min=3"`
	Email                string  `json:"email" validate:"required,email"`
	Password             string  `json:"password" validate:"omitempty,min=8"`
	PasswordConfirmation string  `json:"password_confirmation" validate:"omitempty,eqfield=Password"`
	Phone                string  `json:"phone" validate:"number"`
	Address              string  `json:"address"`
	Lat                  float64 `json:"lat"`
	Lng                  float64 `json:"lng"`
	Photo                string  `json:"photo"`
	RoleId               int64   `json:"role_id" validate:"required"`
}
