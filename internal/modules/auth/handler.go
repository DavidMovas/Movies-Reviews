package auth

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
	"github.com/labstack/echo"
	"gopkg.in/validator.v2"
)

type Handler struct {
	authService *Service
}

func NewHandler(authService *Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) Register(c echo.Context) error {
	req, err := echox.BindAndValidate[RegisterRequest](c)
	if err != nil {
		return err
	}
	if err := validator.Validate(&req); err != nil {
		return apperrors.BadRequestHidden(err, "invalid email or password")
	}

	user := &users.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     users.UserRole,
	}

	if err := h.authService.Register(c.Request().Context(), user, req.Password); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c echo.Context) error {
	req, err := echox.BindAndValidate[LoginRequest](c)
	if err != nil {
		return err
	}

	if err := validator.Validate(&req); err != nil {
		return apperrors.BadRequestHidden(err, "invalid email or password")
	}

	token, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	return c.JSON(http.StatusOK, LoginResponse{AccessToken: token})
}

type RegisterRequest struct {
	Username string `json:"username" validate:"min=3,max=24"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
