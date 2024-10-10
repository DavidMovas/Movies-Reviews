package auth

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/internal/jwt"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
	"github.com/labstack/echo"
)

func Self(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userId")
		claims := jwt.GetClaims(c)

		if claims.Role == users.AdminRole || claims.Subject == userId {
			return next(c)
		}

		return echo.NewHTTPError(http.StatusForbidden)
	}
}

func Editor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userId")
		claims := jwt.GetClaims(c)

		if claims.Role == users.EditorRole || claims.Role == users.AdminRole || claims.Subject == userId {
			return next(c)
		}

		return echo.NewHTTPError(http.StatusForbidden)
	}
}

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)

		if claims.Role == users.AdminRole {
			return next(c)
		}

		return echo.NewHTTPError(http.StatusForbidden)
	}
}
