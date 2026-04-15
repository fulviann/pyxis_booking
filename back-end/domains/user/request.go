package user

import (
	"mime/multipart"
	"time"

	"github.com/devanadindra/signlink-mobile/back-end/utils/constants"
)

type LoginReq struct {
	Email    string         `json:"email" valdiate:"required"`
	Password string         `json:"password" validate:"required"`
	Role     constants.ROLE `json:"role" validate:"omitempty,oneof=ADMIN CUSTOMER"`
}

type LogoutReq struct {
	Token   string
	Expires time.Time
}

type RegisterReq struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type GoogleAuth struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Picture  string `json:"picture"`
	GoogleID string `json:"google_id" validate:"required"`
}

type ChangePasswordReq struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}

type UpdateProfileReq struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
}

type AvatarReq struct {
	AvatarUrl *multipart.FileHeader `json:"avatar" validate:"required"`
}

type ResetPasswordReq struct {
	Email string         `json:"email" validate:"required"`
	Role  constants.ROLE `json:"role" validate:"omitempty,oneof=ADMIN CUSTOMER"`
}

type ResetPasswordSubmitReq struct {
	Email       string         `json:"email" validate:"required"`
	NewPassword string         `json:"newPassword" validate:"required"`
	Role        constants.ROLE `json:"role" validate:"omitempty,oneof=ADMIN CUSTOMER"`
}
