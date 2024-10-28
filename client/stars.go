package client

import (
	"github.com/DavidMovas/Movies-Reviews/contracts"
)

func (c *Client) GetStars(req *contracts.GetStarsRequest) (*contracts.PaginatedResponse[*contracts.Star], error) {
	var stars *contracts.PaginatedResponse[*contracts.Star]

	_, err := c.client.R().
		SetResult(&stars).
		SetQueryParams(req.ToQueryParams()).
		Get(c.path("/api/stars"))

	return stars, err
}

func (c *Client) GetStarByID(req *contracts.GetStarRequest) (*contracts.Star, error) {
	var star contracts.Star

	_, err := c.client.R().
		SetResult(&star).
		Get(c.path("/api/stars/%d", req.StarID))

	return &star, err
}

func (c *Client) CreateStar(req *contracts.AuthenticatedRequest[*contracts.CreateStarRequest]) (*contracts.Star, error) {
	var star *contracts.Star

	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		SetResult(&star).
		Post(c.path("/api/stars"))

	return star, err
}

func (c *Client) UpdateStarByID(req *contracts.AuthenticatedRequest[*contracts.UpdateStarRequest]) (*contracts.Star, error) {
	var star *contracts.Star

	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		SetResult(&star).
		Put(c.path("/api/stars/%d", req.Request.StarID))

	return star, err
}

func (c *Client) DeleteStarByID(accessToken string, req *contracts.DeleteStarRequest) error {
	_, err := c.client.R().
		SetAuthToken(accessToken).
		Delete(c.path("/api/stars/%d", req.StarID))

	return err
}
