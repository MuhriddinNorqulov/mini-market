package dto

type RegisterInput struct {
	Username        string `json:"username" validate:"required,min=3,max=32"`
	Password        string `json:"password" validate:"required,min=6,max=64"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthOutput struct {
	AccessToken string `json:"access_token"`
}
