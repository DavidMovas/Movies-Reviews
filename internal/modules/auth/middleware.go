package auth

import (
	"github.com/DavidMovas/Movies-Reviews/contracts"
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/DavidMovas/Movies-Reviews/internal/jwt"
	"github.com/labstack/echo"
)

var errForbidden = apperrors.Forbidden("insufficient permissions")

func Self(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("userId")
		claims := jwt.GetClaims(c)

		if claims == nil {
			return errForbidden
		}

		if claims.Role == contracts.AdminRole || claims.Subject == userID {
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
		case contracts.AdminRole, contracts.EditorRole:
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

		if claims.Role == contracts.AdminRole {
			return next(c)
		}

		return errForbidden
	}
}
