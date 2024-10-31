package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DavidMovas/Movies-Reviews/client"
	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
)

var (
	titanicReview1   *contracts.Review
	titanicReview2   *contracts.Review
	starsWarsReview1 *contracts.Review
	starsWarsReview2 *contracts.Review
)

func reviewsAPIChecks(t *testing.T, c *client.Client, _ *config.Config) {
	t.Run("reviews.GetReviewsByMovieID: null", func(t *testing.T) {
		res, err := c.GetReviewsByMovieID(&contracts.GetReviewsByMovieIDRequest{
			MovieID: 0,
		})
		require.NoError(t, err)
		require.Equal(t, 0, res.Total)
		require.Equal(t, 0, len(res.Items))
	})

	t.Run("reviews.GetReviewById: not found", func(t *testing.T) {
		_, err := c.GetReviewByID(&contracts.GetReviewRequest{
			ReviewID: 1,
		})
		requireNotFoundError(t, err, "review", "id", 1)
	})

	t.Run("reviews.CreateReview: insufficient permissions", func(t *testing.T) {
		req := &contracts.CreateReviewRequest{
			MovieID: 1,
			UserID:  1,
			Title:   "title",
			Content: "content",
			Rating:  5,
		}

		_, err := c.CreateReview(*contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("reviews.CreateReview: 4 reviews: success", func(t *testing.T) {
		cases := []struct {
			token string
			req   *contracts.CreateReviewRequest
			addr  **contracts.Review
		}{
			{
				token: johnMooreToken,
				req: &contracts.CreateReviewRequest{
					MovieID: titanic.ID,
					UserID:  johnMoore.ID,
					Title:   "Titanic Review",
					Content: "Some content",
					Rating:  8,
				},
				addr: &titanicReview1,
			},
			{
				token: markTwainToken,
				req: &contracts.CreateReviewRequest{
					MovieID: titanic.ID,
					UserID:  markTwain.ID,
					Title:   "Titanic Review",
					Content: "Some content",
					Rating:  2,
				},
				addr: &titanicReview2,
			},

			{
				token: johnMooreToken,
				req: &contracts.CreateReviewRequest{
					MovieID: starWars.ID,
					UserID:  johnMoore.ID,
					Title:   "Stars Wars Review",
					Content: "Some content",
					Rating:  5,
				},
				addr: &starsWarsReview1,
			},
			{
				token: markTwainToken,
				req: &contracts.CreateReviewRequest{
					MovieID: starWars.ID,
					UserID:  markTwain.ID,
					Title:   "Stars Wars Review",
					Content: "Some content",
					Rating:  9,
				},
				addr: &starsWarsReview2,
			},
		}

		for _, cc := range cases {
			res, err := c.CreateReview(*contracts.NewAuthenticated(cc.req, cc.token))
			require.NoError(t, err)
			require.NotNil(t, res)
			require.Equal(t, cc.req.Title, res.Title)
			require.Equal(t, cc.req.Content, res.Content)
			require.Equal(t, cc.req.Rating, res.Rating)
			require.NotNil(t, res.CreatedAt)
			*cc.addr = res
		}
	})

	t.Run("reviews.GetReviewsByMovieID: success", func(t *testing.T) {
		req := &contracts.GetReviewsByMovieIDRequest{
			MovieID: titanic.ID,
		}
		res, err := c.GetReviewsByMovieID(req)
		require.NoError(t, err)
		require.Equal(t, 2, len(res.Items))
		require.Equal(t, 2, res.Total)

		require.Equal(t, titanicReview2.Title, res.Items[0].Title)
		require.Equal(t, titanicReview2.Title, res.Items[1].Title)

		require.Equal(t, titanicReview1.Content, res.Items[0].Content)
		require.Equal(t, titanicReview1.Content, res.Items[1].Content)
	})

	t.Run("reviews.GetReviewsByMovieID: sort by rating (DESC): success", func(t *testing.T) {
		req := &contracts.GetReviewsByMovieIDRequest{
			MovieID: starWars.ID,
			PaginatedRequestOrdered: contracts.PaginatedRequestOrdered{
				Sort:  "rating",
				Order: "desc",
			},
		}
		res, err := c.GetReviewsByMovieID(req)
		require.NoError(t, err)
		require.Equal(t, 2, len(res.Items))
		require.Equal(t, 2, res.Total)

		require.True(t, res.Items[0].Rating > res.Items[1].Rating)

		require.Equal(t, starsWarsReview2.Title, res.Items[0].Title)
		require.Equal(t, starsWarsReview1.Title, res.Items[1].Title)

		require.Equal(t, starsWarsReview2.Content, res.Items[0].Content)
		require.Equal(t, starsWarsReview1.Content, res.Items[1].Content)
	})

	t.Run("reviews.GetReviewById: success", func(t *testing.T) {
		req := &contracts.GetReviewRequest{
			ReviewID: starsWarsReview1.ID,
		}
		res, err := c.GetReviewByID(req)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, starsWarsReview1.Title, res.Title)
		require.Equal(t, starsWarsReview1.Content, res.Content)
		require.Equal(t, starsWarsReview1.Rating, res.Rating)

		req = &contracts.GetReviewRequest{
			ReviewID: titanicReview2.ID,
		}
		res, err = c.GetReviewByID(req)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, titanicReview2.Title, res.Title)
		require.Equal(t, titanicReview2.Content, res.Content)
		require.Equal(t, titanicReview2.Rating, res.Rating)
	})

	t.Run("reviews.GetReviewsByUserID: success", func(t *testing.T) {
		req := &contracts.GetReviewsByUserIDRequest{
			UserID: johnMoore.ID,
		}

		res, err := c.GetReviewsByUserID(req)
		require.NoError(t, err)
		require.Equal(t, 2, len(res.Items))
		require.Equal(t, 2, res.Total)

		require.Equal(t, starsWarsReview1.Title, res.Items[0].Title)
		require.Equal(t, starsWarsReview1.Content, res.Items[0].Content)
		require.Equal(t, starsWarsReview1.Rating, res.Items[0].Rating)

		require.Equal(t, titanicReview1.Title, res.Items[1].Title)
		require.Equal(t, titanicReview1.Content, res.Items[1].Content)
		require.Equal(t, titanicReview1.Rating, res.Items[1].Rating)
	})

	t.Run("reviews.UpdateReview: insufficient permissions", func(t *testing.T) {
		req := &contracts.UpdateReviewRequest{
			ReviewID: titanicReview1.ID,
			Title:    ptr("title"),
			Content:  ptr("content"),
			Rating:   ptr(5),
		}
		_, err := c.UpdateReview(*contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("reviews.UpdateReview: not found", func(t *testing.T) {
		wrongReviewID := 100
		req := &contracts.UpdateReviewRequest{
			UserID:   johnMoore.ID,
			MovieID:  titanic.ID,
			ReviewID: wrongReviewID,
			Title:    ptr("title"),
			Content:  ptr("content"),
			Rating:   ptr(5),
		}
		_, err := c.UpdateReview(*contracts.NewAuthenticated(req, adminToken))
		requireNotFoundError(t, err, "review", "id", wrongReviewID)
	})

	t.Run("reviews.UpdateReview: success", func(t *testing.T) {
		req := &contracts.UpdateReviewRequest{
			UserID:   johnMoore.ID,
			ReviewID: titanicReview1.ID,
			MovieID:  titanic.ID,
			Title:    ptr("title"),
			Content:  ptr("content"),
			Rating:   ptr(10),
		}
		res, err := c.UpdateReview(*contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, *req.Title, res.Title)
		require.Equal(t, *req.Content, res.Content)
		require.Equal(t, *req.Rating, res.Rating)
	})

	t.Run("reviews.DeleteReview: insufficient permissions", func(t *testing.T) {
		req := &contracts.DeleteReviewRequest{
			UserID:   johnMoore.ID,
			ReviewID: titanicReview1.ID,
		}
		err := c.DeleteReview(*contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("reviews.DeleteReview: not found", func(t *testing.T) {
		wrongReviewID := 100
		req := &contracts.DeleteReviewRequest{
			UserID:   johnMoore.ID,
			ReviewID: wrongReviewID,
		}
		err := c.DeleteReview(*contracts.NewAuthenticated(req, johnMooreToken))
		requireNotFoundError(t, err, "review", "id", wrongReviewID)
	})

	t.Run("reviews.DeleteReview: success", func(t *testing.T) {
		req := &contracts.DeleteReviewRequest{
			UserID:   johnMoore.ID,
			ReviewID: titanicReview1.ID,
		}
		err := c.DeleteReview(*contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
	})
}
