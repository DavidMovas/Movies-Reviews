package client

import "github.com/DavidMovas/Movies-Reviews/contracts"

func (c *Client) GetUserById(userId int) (*contracts.UserWithPassword, error) {
	var user *contracts.UserWithPassword

	_, err := c.client.R().
		SetResult(&user).
		Get(c.path("/api/users/%d", userId))

	return user, err
}

func (c *Client) GetUserByUsername(username string) (*contracts.UserWithPassword, error) {
	var user *contracts.UserWithPassword

	_, err := c.client.R().
		SetResult(&user).
		Get(c.path("/api/users/username/%s", username))

	return user, err
}

func (c *Client) UpdateUserData(req *contracts.AuthenticatedRequest[*contracts.UpdateUserRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Request).
		Put(c.path("/api/users/%d", req.Request.UserId))

	return err
}
