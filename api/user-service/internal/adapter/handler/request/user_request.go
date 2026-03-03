package request

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Passowrd string `json:"password" validate:"required,min=8"`
}

type SignUpRequest struct {
	Name                 string `json:"name" validate:"required,min=3"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"password,omitempty"`
	NewPassword     string `json:"password_new" validate:"required,min=8"`
	ConfirmPassword string `json:"password_confirm" validate:"required,min=8"`
}

type UpdateDataRequest struct {
	Name    string  `json:"name" validate:"required,min=3"`
	Email   string  `json:"email" validate:"required,email"`
	Phone   int64   `json:"phone"`
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Photo   string  `json:"photo"`
}

type CustomerBatchRequest struct {
	IDUsers []int64 `json:"id_users" validate:"required"`
}
