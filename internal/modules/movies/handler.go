package movies

import (
	"net/http"

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

func (h *Handler) GetsMovies(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

func (h *Handler) GetMovieByID(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

func (h *Handler) CreateMovie(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

func (h *Handler) UpdateMovieByID(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

func (h *Handler) DeleteMovieByID(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}
