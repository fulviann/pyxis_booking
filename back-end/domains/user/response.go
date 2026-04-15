package user

import (
	"time"

	"github.com/devanadindra/signlink-mobile/back-end/utils/constants"
	"github.com/google/uuid"
)

type LoginRes struct {
	Role    constants.ROLE `json:"role"`
	Token   string         `json:"token"`
	Expires time.Time      `json:"expires"`
}

type VerifyTokenRes struct {
	TokenVerified bool `json:"tokenVerified"`
}

type LogoutRes struct {
	LoggedOut bool `json:"loggedOut"`
}

type ActivityRes struct {
	UserID      uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
}

type PersonalRes struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	AvatarUrl   string    `json:"url"`
	GoogleID    string    `json:"google_id"`
	HasPassword bool      `json:"has_password"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ResetPasswordRes struct {
	Email string
}
