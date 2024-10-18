package client

import (
	"github.com/DavidMovas/Movies-Reviews/contracts"
)

func (c *Client) GetStars() ([]*contracts.Star, error) {
	var stars []*contracts.Star

	_, err := c.client.R().
		SetResult(&stars).
		Get(c.path("/api/stars"))

	return stars, err
}

func (c *Client) GetStarByID(starID int) (*contracts.Star, error) {
	var star contracts.Star

	_, err := c.client.R().
		SetResult(&star).
		Get(c.path("/api/stars/%d", starID))

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

func (c *Client) UpdateStarByID(req *contracts.AuthenticatedRequest[*contracts.UpdateStarRequest]) error {
	var star *contracts.Star

	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		SetResult(&star).
		Put(c.path("/api/stars/%d", req.Request.StarID))

	return err
}

func (c *Client) DeleteStarByID(accessToken string, starID int) error {
	_, err := c.client.R().
		SetAuthToken(accessToken).
		Delete(c.path("/api/stars/%d", starID))

	return err
}
