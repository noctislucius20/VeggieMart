package request

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type SignUpRequest struct {
	Name                 string `json:"name" validate:"required,min=3,max=255"`
	Email                string `json:"email" validate:"required,email,max=255"`
	Password             string `json:"password" validate:"required,min=8,max=255"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,max=255,eqfield=Password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"password" validate:"omitempty,max=255"`
	NewPassword     string `json:"password_new" validate:"required,min=8,max=255"`
	ConfirmPassword string `json:"password_confirm" validate:"required,min=8,max=255,eqfield=NewPassword"`
}

type UpdateDataRequest struct {
	Name    string  `json:"name" validate:"required,min=3,max=255"`
	Email   string  `json:"email" validate:"required,email,max=255"`
	Phone   string  `json:"phone" validate:"omitempty,number,max=17"`
	Address string  `json:"address" validate:"omitempty"`
	Lat     float64 `json:"lat" validate:"omitempty"`
	Lng     float64 `json:"lng" validate:"omitempty"`
	Photo   string  `json:"photo" validate:"omitempty,max=255"`
}

type CustomerBatchRequest struct {
	IDUsers []int64 `json:"id_users" validate:"required,min=1"`
}
