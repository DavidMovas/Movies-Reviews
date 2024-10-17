package integration_tests

import (
	"fmt"
	"testing"

	"github.com/DavidMovas/Movies-Reviews/client"
	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/stretchr/testify/require"
)

func usersApiChecks(t *testing.T, c *client.Client, cfg *config.Config) {
	t.Run("users.GetExistingUserByUsername: not found", func(t *testing.T) {
		_, err := c.GetUserByUsername("someTestName")
		requireNotFoundError(t, err, "user", "username", "someTestName")
	})

	t.Run("users.GetExistingUserById: not found", func(t *testing.T) {
		_, err := c.GetUserById(100)
		requireNotFoundError(t, err, "user", "id", 100)
	})

	t.Run("users.UpdateExistingUserById: success", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserId:   johnMoore.ID,
			Username: fmt.Sprintf("%sTEST", johnMoore.Username),
			Password: fmt.Sprintf("%sTEST", johnMoorePass),
		}
		err := c.UpdateUserData(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
	})

	t.Run("users.UpdateExistingUserById: rollback success", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserId:   johnMoore.ID,
			Username: fmt.Sprintf("%s", johnMoore.Username),
			Password: fmt.Sprintf("%s", johnMoorePass),
		}
		err := c.UpdateUserData(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
	})

	t.Run("users.UpdateExistingUserById: another user", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserId:   johnMoore.ID + 100,
			Username: fmt.Sprintf("%sTEST", johnMoore.Username),
			Password: fmt.Sprintf("%sTEST", johnMoorePass),
		}
		err := c.UpdateUserData(contracts.NewAuthenticated(req, johnMooreToken))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("users.UpdateExistingUserById: non-authenticated", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserId:   johnMoore.ID + 1,
			Username: fmt.Sprintf("%sTEST", johnMoore.Username),
			Password: fmt.Sprintf("%sTEST", johnMoorePass),
		}
		err := c.UpdateUserData(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("users.GetExistingUserByUsername: success", func(t *testing.T) {
		user, err := c.GetUserByUsername(johnMoore.Username)
		require.NoError(t, err)
		require.Equal(t, johnMoore.Username, user.Username)
	})

	t.Run("users.GetExistingUserByUsername: not found", func(t *testing.T) {
		_, err := c.GetUserByUsername("someTestName")
		requireNotFoundError(t, err, "user", "username", "someTestName")
	})

	adminToken = login(t, c, cfg.Admin.Email, cfg.Admin.Password)

	t.Run("users.UpdateUserRoleById: success", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserId:   johnMoore.ID,
			Username: fmt.Sprintf("%s", johnMoore.Username),
			Password: fmt.Sprintf("%s", johnMoorePass),
		}
		err := c.UpdateUserRole(contracts.NewAuthenticated(req, adminToken), contracts.EditorRole)
		require.NoError(t, err)
	})

	johnMooreToken = login(t, c, johnMoore.Email, johnMoorePass)

	t.Run("users.UpdateUserRoleById: invalid role", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserId:   johnMoore.ID,
			Username: fmt.Sprintf("%s", johnMoore.Username),
			Password: fmt.Sprintf("%s", johnMoorePass),
		}
		err := c.UpdateUserRole(contracts.NewAuthenticated(req, adminToken), "invalid_role")
		requireBadRequestError(t, err, "role unknown")
	})

	t.Run("users.UpdateUserRoleById: user not found", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserId:   johnMoore.ID + 100,
			Username: fmt.Sprintf("%s", johnMoore.Username),
			Password: fmt.Sprintf("%s", johnMoorePass),
		}
		err := c.UpdateUserRole(contracts.NewAuthenticated(req, adminToken), contracts.EditorRole)
		requireNotFoundError(t, err, "user", "id", johnMoore.ID+100)
	})

	t.Run("users.UpdateUserRoleById: non-authenticated", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserId:   johnMoore.ID,
			Username: fmt.Sprintf("%s", johnMoore.Username),
			Password: fmt.Sprintf("%s", johnMoorePass),
		}
		err := c.UpdateUserRole(contracts.NewAuthenticated(req, ""), contracts.EditorRole)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("users.DeleteUserById: non-authenticated", func(t *testing.T) {
		err := c.DeleteUserById("", johnMoore.ID)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("users.DeleteUserById: user not found", func(t *testing.T) {
		err := c.DeleteUserById(adminToken, johnMoore.ID+1000)
		requireNotFoundError(t, err, "user", "id", johnMoore.ID+1000)
	})

	t.Run("users.DeleteUserById: success", func(t *testing.T) {
		user := register(t, c, "userForDelete", "2lG8G@example.com", standardPassword)
		err := c.DeleteUserById(adminToken, user.ID)
		require.NoError(t, err)
	})
}
