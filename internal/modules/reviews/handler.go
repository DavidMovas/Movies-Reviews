package reviews

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service          *Service
	paginationConfig *config.PaginationConfig
}

func NewHandler(service *Service, paginationConfig *config.PaginationConfig) *Handler {
	return &Handler{
		service:          service,
		paginationConfig: paginationConfig,
	}
}

func (h *Handler) GetReviewsByMovieID(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

func (h *Handler) GetReviewsByUserID(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

func (h *Handler) GetReviewByID(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

func (h *Handler) CreateReview(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

func (h *Handler) UpdateReview(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

func (h *Handler) DeleteReview(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}
