package dtos

type LoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	AppId    int64  `json:"app_id" validate:"required"`
}

type RegisterDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshDto struct {
	RefreshToken string `json:"refresh_token" validate:"required,jwt"`
	AppId        int64  `json:"app_id" validate:"required"`
}