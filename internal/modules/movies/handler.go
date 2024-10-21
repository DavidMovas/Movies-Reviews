package movies

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	"github.com/DavidMovas/Movies-Reviews/internal/pagination"

	"github.com/labstack/echo"
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

func (h *Handler) GetMovies(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetMoviesRequest](c)
	if err != nil {
		return err
	}

	pagination.SetDefaultsOrdered(&req.PaginatedRequestOrdered, h.paginationConfig)
	offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)

	movies, total, err := h.service.GetMovies(c.Request().Context(), offset, limit, req.Sort, req.Order)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pagination.ResponseOrdered[*contracts.Movie](&req.PaginatedRequestOrdered, total, movies))
}

func (h *Handler) GetMovieByID(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetMovieRequest](c)
	if err != nil {
		return err
	}

	movie, err := h.service.GetMovieByID(c.Request().Context(), req.MovieID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movie)
}

func (h *Handler) CreateMovie(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.CreateMovieRequest](c)
	if err != nil {
		return err
	}

	movie, err := h.service.CreateMovie(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, movie)
}

func (h *Handler) UpdateMovieByID(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateMovieRequest](c)
	if err != nil {
		return err
	}

	movie, err := h.service.UpdateMovieByID(c.Request().Context(), req.MovieID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movie)
}

func (h *Handler) DeleteMovieByID(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.DeleteMovieRequest](c)
	if err != nil {
		return err
	}

	if err = h.service.DeleteMovieByID(c.Request().Context(), req.MovieID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
