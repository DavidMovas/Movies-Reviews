package client

import "github.com/DavidMovas/Movies-Reviews/contracts"

func (c *Client) GetReviewsByMovieID(req *contracts.GetReviewsByMovieIDRequest) (*contracts.PaginatedResponseOrdered[*contracts.Review], error) {
	var resp *contracts.PaginatedResponseOrdered[*contracts.Review]

	_, err := c.client.R().
		SetResult(&resp).
		SetQueryParams(req.ToQueryParams()).
		Get(c.path("/api/movies/%d/reviews", req.MovieID))

	return resp, err
}

func (c *Client) GetReviewsByUserID(req *contracts.GetReviewsByUserIDRequest) (*contracts.PaginatedResponseOrdered[*contracts.Review], error) {
	var resp *contracts.PaginatedResponseOrdered[*contracts.Review]

	_, err := c.client.R().
		SetResult(&resp).
		SetQueryParams(req.ToQueryParams()).
		Get(c.path("/api/users/%d/reviews", req.UserID))

	return resp, err
}

func (c *Client) GetReviewByID(req *contracts.GetReviewRequest) (*contracts.Review, error) {
	var resp *contracts.Review

	_, err := c.client.R().
		SetResult(&resp).
		Get(c.path("/api/reviews/%d", req.ReviewID))

	return resp, err
}

func (c *Client) CreateReview(req contracts.AuthenticatedRequest[*contracts.CreateReviewRequest]) (*contracts.Review, error) {
	var resp *contracts.Review

	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		SetResult(&resp).
		Post(c.path("/api/users/%d/reviews", req.Request.UserID))

	return resp, err
}

func (c *Client) UpdateReview(req contracts.AuthenticatedRequest[*contracts.UpdateReviewRequest]) (*contracts.Review, error) {
	var resp *contracts.Review

	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		SetResult(&resp).
		Put(c.path("/api/users/%d/reviews/%d", req.Request.UserID, req.Request.ReviewID))

	return resp, err
}

func (c *Client) DeleteReview(req contracts.AuthenticatedRequest[*contracts.DeleteReviewRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		Delete(c.path("/api/users/%d/reviews/%d", req.Request.UserID, req.Request.ReviewID))

	return err
}
