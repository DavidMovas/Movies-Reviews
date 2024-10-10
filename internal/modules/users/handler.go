package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"gopkg.in/validator.v2"
)

var (
	ErrInvalidUserId = echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	ErrInvalidRole   = echo.NewHTTPError(http.StatusBadRequest, "invalid role")
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetExistingUsers(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

func (h *Handler) GetExistingUserById(c echo.Context) error {
	userId, err := readUserId(c)
	if err != nil {
		return ErrInvalidUserId
	}

	user, err := h.service.GetExistingUserById(c.Request().Context(), userId)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) GetExistingUserByUsername(c echo.Context) error {
	username := c.Param("username")

	user, err := h.service.GetExistingUserByUsername(c.Request().Context(), username)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateExistingUserById(c echo.Context) error {
	userId, err := readUserId(c)
	if err != nil {
		return ErrInvalidUserId
	}

	isUserExists, err := h.service.CheckUserExistsById(c.Request().Context(), userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	} else if !isUserExists {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	var ud NewUserData
	if err := c.Bind(&ud); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := validator.Validate(&ud); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return h.service.UpdateExistingUserById(c.Request().Context(), userId, &ud)
}

func (h *Handler) UpdateUserRoleById(c echo.Context) error {
	userId, err := readUserId(c)
	if err != nil {
		return ErrInvalidUserId
	}

	newRole := c.Param("role")

	if !ValidateRole(newRole) {
		return ErrInvalidRole
	}

	return h.service.UpdateUserRoleById(c.Request().Context(), userId, newRole)
}

func (h *Handler) DeleteExistingUserById(c echo.Context) error {
	userId, err := readUserId(c)
	if err != nil {
		return err
	}

	return h.service.DeleteExistingUserById(c.Request().Context(), userId)
}

func readUserId(c echo.Context) (int, error) {
	userId := c.Param("userId")
	if userId == "" {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	return id, nil
}

type NewUserData struct {
	Username string `json:"username" validate:"min=3,max=24"`
	Password string `json:"password" validate:"password"`
}
