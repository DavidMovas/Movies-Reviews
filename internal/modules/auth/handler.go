package auth

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
	"github.com/labstack/echo"
)

type Handler struct {
	authService *Service
}

func (h *Handler) Register(c echo.Context) error {
	var raq RegisterRequest
	if err := c.Bind(&raq); err != nil {
		return err
	}

	user := &users.User{
		Username: raq.Username,
		Email:    raq.Email,
	}

	if err := h.authService.Register(c.Request().Context(), user, raq.Password); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c echo.Context) error {
	return c.String(http.StatusOK, "[LOGIN] not implemented")
}

func NewHandler(authService *Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
