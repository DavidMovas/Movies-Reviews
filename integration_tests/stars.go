package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/config"

	"github.com/DavidMovas/Movies-Reviews/client"
)

var (
	jackStar   *contracts.Star
	denzelStar *contracts.Star
	sophiaStar *contracts.Star
)

func starsAPIChecks(t *testing.T, client *client.Client, _ *config.Config) {
	t.Run("stars.GetStars: empty: success", func(t *testing.T) {
		stars, err := client.GetStars()
		require.NoError(t, err)
		require.Nil(t, stars)
	})

	t.Run("stars.GetStarById: not found", func(t *testing.T) {
		_, err := client.GetStarByID(1)
		requireNotFoundError(t, err, "star", "id", 1)
	})

	t.Run("stars.CreateStar: insufficient permissions", func(t *testing.T) {
		req := &contracts.CreateStarRequest{
			FirstName: "Test",
			LastName:  "Test",
		}
		_, err := client.CreateStar(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("stars.CreateStar: create 3 stars: success", func(t *testing.T) {
		cases := []struct {
			req  *contracts.CreateStarRequest
			addr **contracts.Star
		}{
			{
				req: &contracts.CreateStarRequest{
					FirstName:  "jack",
					LastName:   "nicholson",
					BirthDate:  time.Date(1937, 04, 22, 0, 0, 0, 0, time.UTC),
					BirthPlace: ptr("Neptune, New Jersey, USA"),
					Bio:        ptr("Jack Nicholson, an American actor, producer, director and screenwriter, is a three-time Academy Award winner and twelve-time nominee. Nicholson is also notable for being one of two actors - the other being Michael Caine - who have received an Oscar nomination in every decade from the '60s through the '00s."),
				},
				addr: &jackStar,
			},
			{
				req: &contracts.CreateStarRequest{
					FirstName:  "denzel",
					LastName:   "washington",
					BirthDate:  time.Date(1954, 11, 28, 0, 0, 0, 0, time.UTC),
					BirthPlace: ptr("Mount Vernon, New York, USA"),
					Bio:        ptr("Denzel Hayes Washington, Jr. was born on December 28, 1954 in Mount Vernon, New York. He is the middle of three children of a beautician mother, Lennis, from Georgia, and a Pentecostal minister father, Denzel Washington, Sr., from Virginia. After graduating from high school, Denzel enrolled at Fordham University, intent on a career in journalism."),
				},
				addr: &denzelStar,
			},
			{
				req: &contracts.CreateStarRequest{
					FirstName:  "sophia",
					LastName:   "loren",
					BirthDate:  time.Date(1934, 10, 20, 0, 0, 0, 0, time.UTC),
					BirthPlace: ptr("Rome, Lazio, Italy"),
					Bio:        ptr("Sophia Loren was born as Sofia Scicolone at the Clinica Regina Margherita in Rome on September 20, 1934. Her father Riccardo was married to another woman and refused to marry her mother Romilda Villani, despite the fact that she was the mother of his two children (Sophia and her younger sister Maria Scicolone)."),
				},
				addr: &sophiaStar,
			},
		}

		for _, c := range cases {
			star, err := client.CreateStar(contracts.NewAuthenticated(c.req, johnMooreToken))
			require.NoError(t, err)

			*c.addr = star
			require.NotEmpty(t, star.ID)
			require.NotEmpty(t, star.CreatedAt)
		}
	})

	t.Run("stars.CreateStar: firstname is empty", func(t *testing.T) {
		req := &contracts.CreateStarRequest{
			LastName:  "some",
			BirthDate: time.Date(1970, 04, 22, 0, 0, 0, 0, time.UTC),
		}

		_, err := client.CreateStar(contracts.NewAuthenticated(req, johnMooreToken))
		requireBadRequestError(t, err, "FirstName: less than min")
	})

	t.Run("stars.GetStars: not empty: success", func(t *testing.T) {
		stars, err := client.GetStars()
		require.NoError(t, err)
		require.Equal(t, len([]*contracts.Star{jackStar, denzelStar, sophiaStar}), len(stars))
	})

	t.Run("stars.GetStarById: success", func(t *testing.T) {
		star, err := client.GetStarByID(denzelStar.ID)
		require.NoError(t, err)
		require.Equal(t, denzelStar.ID, star.ID)
	})

	t.Run("stars.UpdateStar: insufficient permissions", func(t *testing.T) {
		req := &contracts.UpdateStarRequest{
			StarID:    jackStar.ID,
			FirstName: ptr(jackStar.FirstName + " updated"),
			LastName:  ptr(jackStar.LastName + " updated"),
		}
		_, err := client.UpdateStarByID(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("stars.UpdateStar: not found", func(t *testing.T) {
		const starID = 1000
		req := &contracts.UpdateStarRequest{
			StarID:    starID,
			FirstName: ptr(sophiaStar.FirstName + " updated"),
			LastName:  ptr(sophiaStar.LastName + " updated"),
		}
		_, err := client.UpdateStarByID(contracts.NewAuthenticated(req, johnMooreToken))
		requireNotFoundError(t, err, "star", "id", starID)
	})

	t.Run("stars.UpdateStar: bio update: success", func(t *testing.T) {
		const bio = "updated bio"
		req := &contracts.UpdateStarRequest{
			StarID: jackStar.ID,
			Bio:    ptr(bio),
		}
		star, err := client.UpdateStarByID(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
		require.Equal(t, bio, *star.Bio)
	})

	t.Run("stars.UpdateStar: names update: success", func(t *testing.T) {
		var (
			firstname  = sophiaStar.FirstName + " updated"
			middleName = "new middle name"
			lastname   = sophiaStar.LastName + " updated"
		)
		req := &contracts.UpdateStarRequest{
			StarID:     sophiaStar.ID,
			FirstName:  ptr(firstname),
			MiddleName: ptr(middleName),
			LastName:   ptr(lastname),
		}
		star, err := client.UpdateStarByID(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
		require.Equal(t, sophiaStar.ID, star.ID)
		require.NotNil(t, star.MiddleName)
		require.Equal(t, firstname, star.FirstName)
		require.Equal(t, middleName, *star.MiddleName)
		require.Equal(t, lastname, star.LastName)
	})

	t.Run("stars.UpdateStar: bio update on empty: success", func(t *testing.T) {
		req := &contracts.UpdateStarRequest{
			StarID: jackStar.ID,
			Bio:    ptr(""),
		}
		star, err := client.UpdateStarByID(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
		require.Nil(t, star.Bio)
	})

	t.Run("stars.DeleteStar: insufficient permissions", func(t *testing.T) {
		err := client.DeleteStarByID("", denzelStar.ID)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("stars.DeleteStar: not found", func(t *testing.T) {
		const starID = 1000
		err := client.DeleteStarByID(johnMooreToken, starID)
		requireNotFoundError(t, err, "star", "id", starID)
	})

	t.Run("stars.DeleteStar: success", func(t *testing.T) {
		err := client.DeleteStarByID(johnMooreToken, denzelStar.ID)
		require.NoError(t, err)
	})
}