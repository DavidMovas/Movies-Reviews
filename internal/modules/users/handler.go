package users

import (
	"errors"
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/labstack/echo"
)

const (
	paramUserID     = "userId"
	paramUsername   = "username"
	invalidUserID   = "invalid userid"
	invalidUsername = "invalid username"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetExistingUserByID @Summary Get existing user by id
// @Description Get existing user by id
// @ID get-existing-user-by-id
// @Tags users
// @Param userId path int true "User ID"
// @Produce json
// @Success 200 {object} contracts.User "User"
// @Failure 400 {object} apperrors.Error "Invalid user id, invalid parameter or missing parameter"
// @Failure 404 {object} apperrors.Error "User not found"
// @Failure 500 {object} apperrors.Error "Internal server error"
// @Router /users/{userId} [get]
func (h *Handler) GetExistingUserByID(c echo.Context) error {
	userID, err := echox.ReadFromParam[int](c, paramUserID, invalidUserID)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	user, err := h.service.GetExistingUserByID(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// GetExistingUserByUsername @Summary Get existing user by username
// @Description Get existing user by username
// @ID get-existing-user-by-username
// @Tags users
// @Param username path string true "Username"
// @Produce json
// @Success 200 {object} contracts.User "User"
// @Failure 400 {object} apperrors.Error "Invalid username, invalid parameter or missing parameter"
// @Failure 404 {object} apperrors.Error "User not found"
// @Failure 500 {object} apperrors.Error "Internal server error"
// @Router /users/{username} [get]
func (h *Handler) GetExistingUserByUsername(c echo.Context) error {
	username, err := echox.ReadFromParam[string](c, paramUsername, invalidUsername)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	user, err := h.service.GetExistingUserByUsername(c.Request().Context(), username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateExistingUserByID @Summary Update existing user by id
// @Description Update existing user by id
// @ID update-existing-user-by-id
// @Tags users
// @Param userId path int true "User ID"
// @Param user body contracts.UpdateUserRequest true "User"
// @Accept json
// @Produce json
// @Success 200 "User updated"
// @Failure 400 {object} apperrors.Error "Invalid user id, invalid parameter or missing parameter"
// @Failure 403 {object} apperrors.Error "Insufficient permissions"
// @Failure 404 {object} apperrors.Error "User not found"
// @Failure 500 {object} apperrors.Error "Internal server error"
// @Router /users/{userId} [put]
func (h *Handler) UpdateExistingUserByID(c echo.Context) error {
	userID, err := echox.ReadFromParam[int](c, paramUserID, invalidUserID)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	raq, err := echox.BindAndValidate[UpdateUserRequest](c)
	if err != nil {
		return err
	}

	if err = h.service.UpdateExistingUserByID(c.Request().Context(), userID, raq); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// UpdateUserRoleByID @Summary Update user role by id
// @Description Update user role by id
// @ID update-user-role-by-id
// @Tags users
// @Param userId path int true "User ID"
// @Param role path string true "Role"
// @Produce json
// @Success 200 "User role updated"
// @Failure 400 {object} apperrors.Error "Invalid user id, invalid parameter or missing parameter"
// @Failure 403 {object} apperrors.Error "Insufficient permissions"
// @Failure 404 {object} apperrors.Error "User not found"
// @Failure 500 {object} apperrors.Error "Internal server error"
// @Router /users/{userId}/role/{role} [put]
func (h *Handler) UpdateUserRoleByID(c echo.Context) error {
	userID, err := echox.ReadFromParam[int](c, paramUserID, invalidUserID)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	newRole := c.Param("role")

	if !contracts.ValidateRole(newRole) {
		return apperrors.BadRequestHidden(errors.New("invalid role"), "role unknown")
	}

	if err := h.service.UpdateUserRoleByID(c.Request().Context(), userID, newRole); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// DeleteExistingUserByID @Summary Delete existing user by id
// @Description Delete existing user by id
// @ID delete-existing-user-by-id
// @Tags users
// @Param userId path int true "User ID"
// @Produce json
// @Success 200 "User deleted (softly deleting)"
// @Failure 400 {object} apperrors.Error "Invalid user id, invalid parameter or missing parameter"
// @Failure 403 {object} apperrors.Error "Insufficient permissions"
// @Failure 404 {object} apperrors.Error "User not found"
// @Failure 500 {object} apperrors.Error "Internal server error"
// @Router /users/{userId} [delete]
func (h *Handler) DeleteExistingUserByID(c echo.Context) error {
	userID, err := echox.ReadFromParam[int](c, paramUserID, invalidUserID)
	if err != nil {
		return err
	}

	return h.service.DeleteExistingUserByID(c.Request().Context(), userID)
}
