package genres

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/labstack/echo"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetGenres(c echo.Context) error {
	genres, err := h.Service.GetGenres(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, genres)
}

func (h *Handler) GetGenreById(c echo.Context) error {
	genreId, err := echox.ReadFromParam[int](c, "genreId", "invalid genreId")
	if err != nil {
		return apperrors.BadRequest(err)
	}

	genre, err := h.Service.GetGenreById(c.Request().Context(), genreId)
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

func (h *Handler) UpdateGenreById(c echo.Context) error {
	genreId, err := echox.ReadFromParam[int](c, "genreId", "invalid genreId")
	if err != nil {
		return apperrors.BadRequest(err)
	}

	raq, err := echox.BindAndValidate[contracts.UpdateGenreRequest](c)
	if err != nil {
		return err
	}

	if err := h.Service.UpdateGenreById(c.Request().Context(), genreId, raq); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) DeleteGenreById(c echo.Context) error {
	genreId, err := echox.ReadFromParam[int](c, "genreId", "invalid genreId")
	if err != nil {
		return apperrors.BadRequest(err)
	}

	return h.Service.DeleteGenreById(c.Request().Context(), genreId)
}
