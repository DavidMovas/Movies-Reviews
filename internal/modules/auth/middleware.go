package auth

import (
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/DavidMovas/Movies-Reviews/internal/jwt"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
	"github.com/labstack/echo"
)

var errForbidden = apperrors.Forbidden("insufficient permissions")

func Self(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userId")
		claims := jwt.GetClaims(c)

		if claims.Role == users.AdminRole || claims.Subject == userId {
			return next(c)
		}

		return errForbidden
	}
}

func Editor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		switch jwt.GetClaims(c).Role {
		case users.AdminRole, users.EditorRole:
			return next(c)
		default:
			return errForbidden
		}
	}
}

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if jwt.GetClaims(c).Role == users.AdminRole {
			return next(c)
		}

		return errForbidden
	}
}
