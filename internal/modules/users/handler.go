package users

import (
	"net/http"

	"github.com/golang/groupcache/singleflight"

	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service

	reqGroup singleflight.Group
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
	res, err := h.reqGroup.Do(c.Request().RequestURI, func() (any, error) {
		req, err := echox.BindAndValidate[GetUserByIDRequest](c)
		if err != nil {
			return nil, err
		}

		user, err := h.service.GetExistingUserByID(c.Request().Context(), req.UserID)
		if err != nil {
			return nil, err
		}

		return user, nil
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
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
	res, err := h.reqGroup.Do(c.Request().RequestURI, func() (any, error) {
		req, err := echox.BindAndValidate[GetUserByUsernameRequest](c)
		if err != nil {
			return nil, err
		}

		user, err := h.service.GetExistingUserByUsername(c.Request().Context(), req.Username)
		if err != nil {
			return nil, err
		}

		return user.User, nil
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
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
	req, err := echox.BindAndValidate[UpdateUserRequest](c)
	if err != nil {
		return err
	}

	user, err := h.service.UpdateExistingUserByID(c.Request().Context(), req.UserID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
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
	req, err := echox.BindAndValidate[UpdateUserRoleRequest](c)
	if err != nil {
		return err
	}

	if err = h.service.UpdateUserRoleByID(c.Request().Context(), req.UserID, req.Role); err != nil {
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
	req, err := echox.BindAndValidate[GetUserByIDRequest](c)
	if err != nil {
		return err
	}

	return h.service.DeleteExistingUserByID(c.Request().Context(), req.UserID)
}
