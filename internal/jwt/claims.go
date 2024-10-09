package jwt

import (
	"time"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
	"github.com/golang-jwt/jwt/v4"
)

type AccessClaims struct {
	jwt.StandardClaims
	UserID uint64 `json:"user_id"`
	Role   string `json:"role"`
}

func NewAccessClaimsFromUser(user *users.User, exp time.Duration) *AccessClaims {
	return &AccessClaims{
		UserID: user.ID,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(exp.Minutes()),
			IssuedAt:  time.Now().Unix(),
			Subject:   user.Email,
		},
	}
}
