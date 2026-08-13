package request

type CustomerRequest struct {
	Name                 string `json:"name" validate:"required,min=3,max=255"`
	Email                string `json:"email" validate:"required,email,max=255"`
	Password             string `json:"password" validate:"required,min=8,max=255"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,max=255,eqfield=Password"`
	Phone                string `json:"phone" validate:"omitempty,number,max=17"`
	Address              string `json:"address" validate:"omitempty"`
	Photo                string `json:"photo" validate:"omitempty,max=255"`
}
