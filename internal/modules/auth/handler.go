package auth

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
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
	req, err := echox.BindAndValidate[contracts.RegisterUserRequest](c)
	if err != nil {
		return err
	}
	if err := validator.Validate(&req); err != nil {
		return apperrors.BadRequestHidden(err, "invalid email or password")
	}

	user := &contracts.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     contracts.UserRole,
	}

	if err := h.authService.Register(c.Request().Context(), user, req.Password); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.LoginUserRequest](c)
	if err != nil {
		return err
	}

	if err := validator.Validate(&req); err != nil {
		return apperrors.BadRequestHidden(err, "invalid email or password")
	}

	token, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contracts.LoginUserResponse{AccessToken: token})
}
