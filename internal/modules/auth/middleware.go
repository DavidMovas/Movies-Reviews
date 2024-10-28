package auth

import (
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/DavidMovas/Movies-Reviews/internal/jwt"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
	"github.com/labstack/echo/v4"
)

var errForbidden = apperrors.Forbidden("insufficient permissions")

func Self(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("userId")
		claims := jwt.GetClaims(c)

		if claims == nil {
			return errForbidden
		}

		if claims.Role == users.AdminRole || claims.Subject == userID {
			return next(c)
		}

		return errForbidden
	}
}

func Editor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)

		if claims == nil {
			return errForbidden
		}

		switch claims.Role {
		case users.AdminRole, users.EditorRole:
			return next(c)
		default:
			return errForbidden
		}
	}
}

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)

		if claims == nil {
			return errForbidden
		}

		if claims.Role == users.AdminRole {
			return next(c)
		}

		return errForbidden
	}
}
