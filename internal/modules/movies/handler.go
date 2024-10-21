package movies

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	"github.com/DavidMovas/Movies-Reviews/internal/pagination"

	"github.com/labstack/echo"
)

const (
	paramMovieID   = "movieId"
	invalidMovieID = "invalid movie id"
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
	if err = contracts.ValidateSortRequest(req.Sort); err != nil {
		req.Sort = "id"
	}

	movies, total, err := h.service.GetMovies(c.Request().Context(), offset, limit, req.Sort, req.Order)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pagination.ResponseOrdered[*contracts.Movie](&req.PaginatedRequestOrdered, total, movies))
}

func (h *Handler) GetMovieByID(c echo.Context) error {
	movieID, err := echox.ReadFromParam[int](c, paramMovieID, invalidMovieID)
	if err != nil {
		return err
	}

	movie, err := h.service.GetMovieByID(c.Request().Context(), movieID)
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
	movieID, err := echox.ReadFromParam[int](c, paramMovieID, invalidMovieID)
	if err != nil {
		return err
	}

	req, err := echox.BindAndValidate[contracts.UpdateMovieRequest](c)
	if err != nil {
		return err
	}

	movie, err := h.service.UpdateMovieByID(c.Request().Context(), movieID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movie)
}

func (h *Handler) DeleteMovieByID(c echo.Context) error {
	movieID, err := echox.ReadFromParam[int](c, paramMovieID, invalidMovieID)
	if err != nil {
		return err
	}

	if err = h.service.DeleteMovieByID(c.Request().Context(), movieID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
