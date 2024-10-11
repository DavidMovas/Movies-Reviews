package jwt

import (
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

const (
	tokenContextKey = "token"
)

func NewAuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenStr := c.Request().Header.Get("Authorization")
			token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil && tokenStr == "" {
				return next(c)
			}

			if err != nil || !token.Valid {
				return apperrors.Forbidden("invalid token")
			}

			c.Set(tokenContextKey, token)

			return next(c)
		}
	}
}

func GetClaims(c echo.Context) *AccessClaims {
	token := c.Get(tokenContextKey)

	if token == nil {
		return nil
	}

	return token.(*jwt.Token).Claims.(*AccessClaims)
}
