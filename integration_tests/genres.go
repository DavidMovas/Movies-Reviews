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
		_, err := c.GetGenreByID(1)
		requireNotFoundError(t, err, "genre", "id", 1)
	})

	t.Run("genres.CreateGenre: insufficient permissions", func(t *testing.T) {
		raq := contracts.CreateGenreRequest{
			Name: "action",
		}

		_, err := c.CreateGenre("", &raq)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.CreateGenre: success", func(t *testing.T) {
		var err error
		raq := contracts.CreateGenreRequest{
			Name: "action",
		}
		actionGenre, err = c.CreateGenre(johnMooreToken, &raq)
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
			genre, err := c.CreateGenre(johnMooreToken, &cc.req)
			require.NoError(t, err)
			require.NotNil(t, cc.addr)
			*cc.addr = genre
			require.NotEmpty(t, genre.ID)
			require.NotEmpty(t, genre.Name)
		}
	})

	t.Run("genres.CreateGenre: already exists", func(t *testing.T) {
		raq := contracts.CreateGenreRequest{
			Name: "action",
		}
		_, err := c.CreateGenre(johnMooreToken, &raq)
		requireAlreadyExistsError(t, err, "genre", "name", "action")
	})

	t.Run("genres.CreateGenre: bad request", func(t *testing.T) {
		raq := contracts.CreateGenreRequest{
			Name: "a",
		}
		_, err := c.CreateGenre(adminToken, &raq)
		requireBadRequestError(t, err, "Name: less than min")
	})

	t.Run("genres.CreateGenre: insufficient permissions", func(t *testing.T) {
		raq := contracts.CreateGenreRequest{
			Name: "action",
		}
		_, err := c.CreateGenre("", &raq)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.GetGenreById: success", func(t *testing.T) {
		requestedGenre, err := c.GetGenreByID(actionGenre.ID)
		require.NoError(t, err)
		require.Equal(t, actionGenre.ID, requestedGenre.ID)
		require.NotNil(t, requestedGenre)
	})

	t.Run("genres.UpdateGenreById: insufficient permissions", func(t *testing.T) {
		raq := contracts.UpdateGenreRequest{
			Name: "horror",
		}
		err := c.UpdateGenreByID("", &raq, actionGenre.ID)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.UpdateGenreById: not found", func(t *testing.T) {
		raq := contracts.UpdateGenreRequest{
			Name: "horror",
		}
		err := c.UpdateGenreByID(johnMooreToken, &raq, 100)
		requireNotFoundError(t, err, "genre", "id", 100)
	})

	t.Run("genres.UpdateGenreById: already exists", func(t *testing.T) {
		raq := contracts.UpdateGenreRequest{
			Name: actionGenre.Name,
		}
		err := c.UpdateGenreByID(johnMooreToken, &raq, actionGenre.ID+1)
		requireAlreadyExistsError(t, err, "genre", "name", actionGenre.Name)
	})

	t.Run("genres.UpdateGenreById: success", func(t *testing.T) {
		var err error
		raq := contracts.UpdateGenreRequest{
			Name: "romance 2",
		}
		err = c.UpdateGenreByID(johnMooreToken, &raq, romanceGenre.ID)
		require.NoError(t, err)
		romanceGenre.Name = raq.Name
	})

	t.Run("genres.DeleteGenreById: insufficient permissions", func(t *testing.T) {
		err := c.DeleteGenreByID("", romanceGenre.ID)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.DeleteGenreById: not found", func(t *testing.T) {
		err := c.DeleteGenreByID(johnMooreToken, romanceGenre.ID+100)
		requireNotFoundError(t, err, "genre", "id", romanceGenre.ID+100)
	})

	t.Run("genres.DeleteGenreById: success", func(t *testing.T) {
		err := c.DeleteGenreByID(johnMooreToken, romanceGenre.ID)
		require.NoError(t, err)
	})
}
