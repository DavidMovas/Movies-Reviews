package jwt

import "github.com/golang-jwt/jwt/v4"

type AccessClaims struct {
	jwt.StandardClaims
	UserID uint64 `json:"user_id"`
	Role   string `json:"role"`
}
