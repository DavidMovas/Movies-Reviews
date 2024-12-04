package jwt

import (
	"strings"

	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type contextKey string

const tokenContextKey contextKey = "token"

func NewAuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenStr := c.Request().Header.Get("Authorization")
			token, err := jwt.ParseWithClaims(clearToken(tokenStr), &AccessClaims{}, func(_ *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil && tokenStr == "" {
				return next(c)
			}

			if err != nil || !token.Valid {
				return apperrors.Forbidden("invalid token")
			}

			c.Set(string(tokenContextKey), token)

			return next(c)
		}
	}
}

func GetClaims(c echo.Context) *AccessClaims {
	token := c.Get(string(tokenContextKey))

	if token == nil {
		return nil
	}

	return token.(*jwt.Token).Claims.(*AccessClaims)
}

func clearToken(tokenStr string) string {
	if strings.Contains(tokenStr, "\"") {
		tokenStr = strings.Trim(tokenStr, "\"")
	}

	if strings.Contains(tokenStr, "Bearer") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}

	return tokenStr
}
