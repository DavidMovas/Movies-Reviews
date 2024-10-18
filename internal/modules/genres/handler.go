package genres

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/labstack/echo"
)

const (
	paramGenreID   = "genreId"
	invalidGenreID = "invalid genreId"
)

type Handler struct {
	*Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetGenres(c echo.Context) error {
	genres, err := h.Service.GetGenres(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, genres)
}

func (h *Handler) GetGenreByID(c echo.Context) error {
	genreID, err := echox.ReadFromParam[int](c, paramGenreID, invalidGenreID)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	genre, err := h.Service.GetGenreByID(c.Request().Context(), genreID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, genre)
}

func (h *Handler) CreateGenre(c echo.Context) error {
	raq, err := echox.BindAndValidate[contracts.CreateGenreRequest](c)
	if err != nil {
		return err
	}

	genre, err := h.Service.CreateGenre(c.Request().Context(), raq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, genre)
}

func (h *Handler) UpdateGenreByID(c echo.Context) error {
	genreID, err := echox.ReadFromParam[int](c, paramGenreID, invalidGenreID)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	raq, err := echox.BindAndValidate[contracts.UpdateGenreRequest](c)
	if err != nil {
		return err
	}

	if err := h.Service.UpdateGenreByID(c.Request().Context(), genreID, raq); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) DeleteGenreByID(c echo.Context) error {
	genreID, err := echox.ReadFromParam[int](c, paramGenreID, invalidGenreID)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	return h.Service.DeleteGenreByID(c.Request().Context(), genreID)
}
