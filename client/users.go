package client

import "github.com/DavidMovas/Movies-Reviews/contracts"

func (c *Client) GetUserByID(userID int) (*contracts.UserWithPassword, error) {
	var user *contracts.UserWithPassword

	_, err := c.client.R().
		SetResult(&user).
		Get(c.path("/api/users/%d", userID))

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
		SetBody(req.Request).
		SetHeader("Content-Type", "application/json").
		Put(c.path("/api/users/%d", req.Request.UserID))

	return err
}

func (c *Client) UpdateUserRole(req *contracts.AuthenticatedRequest[*contracts.UpdateUserRequest], role string) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		Put(c.path("/api/users/%d/role/%s", req.Request.UserID, role))

	return err
}

func (c *Client) DeleteUserByID(accessToken string, userID int) error {
	_, err := c.client.R().
		SetAuthToken(accessToken).
		SetHeader("Content-Type", "application/json").
		Delete(c.path("/api/users/%d", userID))

	return err
}
