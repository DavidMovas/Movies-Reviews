package tests

import (
	"testing"

	"github.com/DavidMovas/Movies-Reviews/client"
	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/stretchr/testify/require"
)

var (
	actionGenre  *contracts.Genre
	dramaGenre   *contracts.Genre
	romanceGenre *contracts.Genre
	comedyGenre  *contracts.Genre
)

func genresAPIChecks(t *testing.T, c *client.Client, _ *config.Config) {
	t.Run("genres.GetGenres: nil", func(t *testing.T) {
		genres, err := c.GetGenres()
		require.NoError(t, err)
		require.Nil(t, genres)
	})

	t.Run("genres.GetGenreById: not found", func(t *testing.T) {
		req := &contracts.GetGenreRequest{GenreID: 1}
		_, err := c.GetGenreByID(req)
		requireNotFoundError(t, err, "genre", "id", 1)
	})

	t.Run("genres.CreateGenre: insufficient permissions", func(t *testing.T) {
		req := contracts.CreateGenreRequest{
			Name: "action",
		}

		_, err := c.CreateGenre(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.CreateGenre: success", func(t *testing.T) {
		var err error
		req := contracts.CreateGenreRequest{
			Name: "action",
		}
		actionGenre, err = c.CreateGenre(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
		require.NotNil(t, actionGenre)
	})

	t.Run("genres.CreateGenre: create 3 genres: success", func(t *testing.T) {
		cases := []struct {
			req  contracts.CreateGenreRequest
			addr **contracts.Genre
		}{
			{
				req: contracts.CreateGenreRequest{
					Name: "drama",
				},
				addr: &dramaGenre,
			},
			{
				req: contracts.CreateGenreRequest{
					Name: "romance",
				},
				addr: &romanceGenre,
			},
			{
				req: contracts.CreateGenreRequest{
					Name: "comedy",
				},
				addr: &comedyGenre,
			},
		}

		for _, cc := range cases {
			genre, err := c.CreateGenre(contracts.NewAuthenticated(cc.req, johnMooreToken))
			require.NoError(t, err)
			require.NotNil(t, cc.addr)
			*cc.addr = genre
			require.NotEmpty(t, genre.ID)
			require.NotEmpty(t, genre.Name)
		}
	})

	t.Run("genres.CreateGenre: already exists", func(t *testing.T) {
		req := contracts.CreateGenreRequest{
			Name: "action",
		}
		_, err := c.CreateGenre(contracts.NewAuthenticated(req, johnMooreToken))
		requireAlreadyExistsError(t, err, "genre", "name", "action")
	})

	t.Run("genres.CreateGenre: bad request", func(t *testing.T) {
		req := contracts.CreateGenreRequest{
			Name: "a",
		}
		_, err := c.CreateGenre(contracts.NewAuthenticated(req, adminToken))
		requireBadRequestError(t, err, "Name: less than min")
	})

	t.Run("genres.CreateGenre: insufficient permissions", func(t *testing.T) {
		req := contracts.CreateGenreRequest{
			Name: "action",
		}
		_, err := c.CreateGenre(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.GetGenreById: success", func(t *testing.T) {
		req := &contracts.GetGenreRequest{GenreID: actionGenre.ID}
		requestedGenre, err := c.GetGenreByID(req)
		require.NoError(t, err)
		require.Equal(t, actionGenre.ID, requestedGenre.ID)
		require.NotNil(t, requestedGenre)
	})

	t.Run("genres.UpdateGenreById: insufficient permissions", func(t *testing.T) {
		req := contracts.UpdateGenreRequest{
			GenreID: actionGenre.ID,
			Name:    "horror",
		}
		err := c.UpdateGenreByID(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.UpdateGenreById: not found", func(t *testing.T) {
		req := contracts.UpdateGenreRequest{
			GenreID: 100,
			Name:    "horror",
		}
		err := c.UpdateGenreByID(contracts.NewAuthenticated(req, johnMooreToken))
		requireNotFoundError(t, err, "genre", "id", req.GenreID)
	})

	t.Run("genres.UpdateGenreById: already exists", func(t *testing.T) {
		req := contracts.UpdateGenreRequest{
			GenreID: actionGenre.ID + 1,
			Name:    actionGenre.Name,
		}
		err := c.UpdateGenreByID(contracts.NewAuthenticated(req, johnMooreToken))
		requireAlreadyExistsError(t, err, "genre", "name", actionGenre.Name)
	})

	t.Run("genres.UpdateGenreById: success", func(t *testing.T) {
		req := contracts.UpdateGenreRequest{
			GenreID: romanceGenre.ID,
			Name:    "romance 2",
		}
		err := c.UpdateGenreByID(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
		romanceGenre.Name = req.Name
	})

	t.Run("genres.DeleteGenreById: insufficient permissions", func(t *testing.T) {
		req := contracts.DeleteGenreRequest{GenreID: romanceGenre.ID}
		err := c.DeleteGenreByID(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.DeleteGenreById: not found", func(t *testing.T) {
		req := contracts.DeleteGenreRequest{GenreID: romanceGenre.ID + 100}
		err := c.DeleteGenreByID(contracts.NewAuthenticated(req, johnMooreToken))
		requireNotFoundError(t, err, "genre", "id", req.GenreID)
	})

	t.Run("genres.DeleteGenreById: success", func(t *testing.T) {
		req := contracts.DeleteGenreRequest{GenreID: romanceGenre.ID}
		err := c.DeleteGenreByID(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
	})
}
