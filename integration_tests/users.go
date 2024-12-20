package tests

import (
	"fmt"
	"testing"

	"github.com/DavidMovas/Movies-Reviews/client"
	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/stretchr/testify/require"
)

func usersAPIChecks(t *testing.T, c *client.Client, cfg *config.Config) {
	t.Run("users.GetExistingUserByUsername: not found", func(t *testing.T) {
		req := &contracts.GetUserByUsernameRequest{Username: "someTestName"}
		_, err := c.GetUserByUsername(req)
		requireNotFoundError(t, err, "user", "username", "someTestName")
	})

	t.Run("users.GetExistingUserById: not found", func(t *testing.T) {
		req := &contracts.GetUserByIDRequest{UserID: 100}
		_, err := c.GetUserByID(req)
		requireNotFoundError(t, err, "user", "id", 100)
	})

	t.Run("users.UpdateExistingUserById: success", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserID:   johnMoore.ID,
			Username: ptr(fmt.Sprintf("%sTEST", johnMoore.Username)),
			Password: ptr(fmt.Sprintf("%sTEST", johnMoorePass)),
		}
		user, err := c.UpdateUserData(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
		require.Equal(t, *req.Username, user.Username)
	})

	t.Run("users.UpdateExistingUserById: rollback success", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserID:   johnMoore.ID,
			Username: ptr(johnMoore.Username),
			Password: ptr(johnMoorePass),
		}
		_, err := c.UpdateUserData(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
	})

	t.Run("users.UpdateExistingUserById: another user", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserID:   johnMoore.ID + 100,
			Username: ptr(fmt.Sprintf("%sTEST", johnMoore.Username)),
			Password: ptr(fmt.Sprintf("%sTEST", johnMoorePass)),
		}
		_, err := c.UpdateUserData(contracts.NewAuthenticated(req, johnMooreToken))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("users.UpdateExistingUserById: non-authenticated", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserID:   johnMoore.ID + 1,
			Username: ptr(fmt.Sprintf("%sTEST", johnMoore.Username)),
			Password: ptr(fmt.Sprintf("%sTEST", johnMoorePass)),
		}
		_, err := c.UpdateUserData(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("users.GetExistingUserByUsername: success", func(t *testing.T) {
		req := &contracts.GetUserByUsernameRequest{Username: johnMoore.Username}
		user, err := c.GetUserByUsername(req)
		require.NoError(t, err)
		require.Equal(t, johnMoore.Username, user.Username)
	})

	t.Run("users.GetExistingUserByUsername: not found", func(t *testing.T) {
		req := &contracts.GetUserByUsernameRequest{Username: "someTestName"}
		_, err := c.GetUserByUsername(req)
		requireNotFoundError(t, err, "user", "username", "someTestName")
	})

	adminToken = login(t, c, cfg.Admin.Email, cfg.Admin.Password)

	t.Run("users.UpdateUserRoleById: success", func(t *testing.T) {
		req := &contracts.UpdateUserRoleRequest{
			UserID: johnMoore.ID,
			Role:   contracts.EditorRole,
		}
		err := c.UpdateUserRole(contracts.NewAuthenticated(req, adminToken))
		require.NoError(t, err)
	})

	johnMooreToken = login(t, c, johnMoore.Email, johnMoorePass)
	markTwainToken = login(t, c, markTwain.Email, markTwainPass)

	t.Run("users.UpdateUserRoleById: invalid role", func(t *testing.T) {
		req := &contracts.UpdateUserRoleRequest{
			UserID: johnMoore.ID,
			Role:   "invalid_role",
		}
		err := c.UpdateUserRole(contracts.NewAuthenticated(req, adminToken))
		requireBadRequestError(t, err, "invalid role")
	})

	t.Run("users.UpdateUserRoleById: user not found", func(t *testing.T) {
		req := &contracts.UpdateUserRoleRequest{
			UserID: johnMoore.ID + 100,
			Role:   contracts.EditorRole,
		}
		err := c.UpdateUserRole(contracts.NewAuthenticated(req, adminToken))
		requireNotFoundError(t, err, "user", "id", johnMoore.ID+100)
	})

	t.Run("users.UpdateUserRoleById: non-authenticated", func(t *testing.T) {
		req := &contracts.UpdateUserRoleRequest{
			UserID: johnMoore.ID,
			Role:   contracts.EditorRole,
		}
		err := c.UpdateUserRole(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("users.DeleteUserById: non-authenticated", func(t *testing.T) {
		req := &contracts.DeleteUserRequest{UserID: johnMoore.ID}
		err := c.DeleteUserByID(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("users.DeleteUserById: user not found", func(t *testing.T) {
		req := &contracts.DeleteUserRequest{UserID: johnMoore.ID + 1000}
		err := c.DeleteUserByID(contracts.NewAuthenticated(req, adminToken))
		requireNotFoundError(t, err, "user", "id", johnMoore.ID+1000)
	})

	t.Run("users.DeleteUserById: success", func(t *testing.T) {
		user := register(t, c, "userForDelete", "2lG8G@example.com", standardPassword)
		req := &contracts.DeleteUserRequest{UserID: user.ID}
		err := c.DeleteUserByID(contracts.NewAuthenticated(req, adminToken))
		require.NoError(t, err)
	})
}
