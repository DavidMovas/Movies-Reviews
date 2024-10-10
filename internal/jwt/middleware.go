package jwt

import (
	"net/http"

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

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, "Invalid token")
			}

			c.Set(tokenContextKey, token)

			return next(c)
		}
	}
}

func GetClaims(c echo.Context) *AccessClaims {
	token := c.Get(tokenContextKey)
	if token == nil {
		panic("token not found in context")
	}

	t, ok := token.(*jwt.Token)
	if !ok {
		panic("invalid token type")
	}

	claims, ok := t.Claims.(*AccessClaims)
	if !ok {
		panic("invalid claims type")
	}
	return claims
}
