package client

import "github.com/DavidMovas/Movies-Reviews/contracts"

func (c *Client) GetUserByID(req *contracts.GetUserByIDRequest) (*contracts.User, error) {
	var user *contracts.User

	_, err := c.client.R().
		SetResult(&user).
		Get(c.path("/api/users/%d", req.UserID))

	return user, err
}

func (c *Client) GetUserByUsername(req *contracts.GetUserByUsernameRequest) (*contracts.User, error) {
	var user *contracts.User

	_, err := c.client.R().
		SetResult(&user).
		Get(c.path("/api/users/username/%s", req.Username))

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

func (c *Client) UpdateUserRole(req *contracts.AuthenticatedRequest[*contracts.UpdateUserRoleRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		Put(c.path("/api/users/%d/role/%s", req.Request.UserID, req.Request.Role))

	return err
}

func (c *Client) DeleteUserByID(req *contracts.AuthenticatedRequest[*contracts.DeleteUserRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		Delete(c.path("/api/users/%d", req.Request.UserID))

	return err
}
