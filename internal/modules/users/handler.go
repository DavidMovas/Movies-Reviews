package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
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
		return err
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
