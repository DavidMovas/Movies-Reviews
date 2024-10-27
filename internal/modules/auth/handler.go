package auth

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"

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

// Register @Summary Register a new user
// @Description Register a new user
// @ID register
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterUserRequest true "User"
// @Success 201 {object} contracts.User "User"
// @Failure 400 {object} apperrors.Error "Invalid email or password"
// @Failure 500 {object} apperrors.Error "Internal server error"
// @Router /auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	req, err := echox.BindAndValidate[RegisterUserRequest](c)
	if err != nil {
		return err
	}
	if err = validator.Validate(&req); err != nil {
		return apperrors.BadRequestHidden(err, "invalid email or password")
	}

	user := &users.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     contracts.UserRole,
	}

	if err = h.authService.Register(c.Request().Context(), user, req.Password); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

// Login @Summary Login a user
// @Description Login a user and return an access token
// @ID login
// @Tags auth
// @Accept json
// @Produce json
// @Param user body LoginUserRequest true "User"
// @Success 200 {object} LoginUserResponse "Access token"
// @Failure 400 {object} apperrors.Error "Invalid email or password"
// @Failure 404 {object} apperrors.Error "User not found"
// @Failure 500 {object} apperrors.Error "Internal server error"
// @Router /auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	req, err := echox.BindAndValidate[LoginUserRequest](c)
	if err != nil {
		return err
	}

	if err = validator.Validate(&req); err != nil {
		return apperrors.BadRequestHidden(err, "invalid email or password")
	}

	token, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, LoginUserResponse{AccessToken: token})
}
