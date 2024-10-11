package echox

import (
	"log/slog"

	"github.com/DavidMovas/Movies-Reviews/internal/jwt"
	"github.com/DavidMovas/Movies-Reviews/internal/log"
	"github.com/labstack/echo"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestGroup := slog.Group("request",
			slog.String("method", c.Request().Method),
			slog.String("url", c.Request().URL.String()),
		)
		attrs := []any{requestGroup}

		if claims := jwt.GetClaims(c); claims != nil {
			requesterGroup := slog.Group("requester",
				slog.Int("id", claims.UserID),
				slog.String("role", claims.Role),
			)

			attrs = append(attrs, requesterGroup)
		}

		logger := slog.Default().With(attrs...)
		ctx := log.WithLogger(c.Request().Context(), logger)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
