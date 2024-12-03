package auth

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"

	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	authService *Service
	userService *users.Service
}

func NewHandler(authService *Service, userService *users.Service) *Handler {
	return &Handler{
		authService: authService,
		userService: userService,
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

	user := &users.User{
		Username:  req.Username,
		Email:     req.Email,
		Role:      users.UserRole,
		AvatarURL: users.DefaultAvatarURL,
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
	req, err := echox.BindAndValidateLoginRequest(c)
	if err != nil {
		return err
	}

	user, token, err := h.authService.Login(c.Request().Context(), req.Email, req.Username, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, LoginUserResponse{AccessToken: token, User: *user.User})
}
