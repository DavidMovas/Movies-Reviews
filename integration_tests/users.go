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
}
