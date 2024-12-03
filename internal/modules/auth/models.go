package auth

import (
	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
)

type RegisterUserRequest struct {
	Username string `json:"username" validate:"min=3,max=24"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}

type LoginUserRequest struct {
	Username *string `json:"username,omitempty" validate:"min=3,max=24"`
	Email    *string `json:"email,omitempty" validate:"email"`
	Password string  `json:"password" validate:"password"`
}

type LoginUserResponse struct {
	User        users.User `json:"user"`
	AccessToken string     `json:"access_token" nonzero:"true"`
}

type RefreshTokenRequest struct {
	UserID      int    `json:"-" param:"userId" validate:"nonzero"`
	AccessToken string `json:"access-token" validate:"nonzero"`
}

type RefreshTokenResponse struct {
	RefreshToken string `json:"access-token" nonzero:"true"`
}

type AuthenticatedRequest[T any] struct {
	AccessToken string
	Request     T
}
