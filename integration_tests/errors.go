package integration_tests

import (
	"errors"
	"net/http"
	"testing"

	"github.com/DavidMovas/Movies-Reviews/client"
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/stretchr/testify/require"
)

func requireNotFoundError(t *testing.T, err error, subject, key string, value any) {
	msg := apperrors.NotFound(subject, key, value).Error()
	requireApiError(t, err, http.StatusNotFound, msg)
}

func requireUnauthorizedError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusUnauthorized, msg)
}

func requireForbiddenError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusForbidden, msg)
}

func requireBadRequestError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusBadRequest, msg)
}

func requireAlreadyExistsError(t *testing.T, err error, subject, key string, value any) {
	msg := apperrors.AlreadyExists(subject, key, value).Error()
	requireApiError(t, err, http.StatusConflict, msg)
}

func requireApiError(t *testing.T, err error, statusCode int, msg string) {
	var cerr *client.Error
	ok := errors.As(err, &cerr)
	require.True(t, ok, "expected client.Error")
	require.Equal(t, statusCode, cerr.Code)
	require.Contains(t, cerr.Message, msg)
}
