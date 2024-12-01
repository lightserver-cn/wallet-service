package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func TestAssertResponse(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("normal response", func(t *testing.T) {
		raw := map[string]any{
			"username": "testuser",
			"email":    "testuser@example.com",
		}
		var res User
		err := AssertResponse(raw, &res)
		require.NoError(t, err)
		assert.Equal(t, "testuser", res.Username)
		assert.Equal(t, "testuser@example.com", res.Email)
	})

	t.Run("response with error field", func(t *testing.T) {
		raw := map[string]any{
			"error": "some error message",
		}
		var res User
		err := AssertResponse(raw, &res)
		require.Error(t, err)
		assert.EqualError(t, err, "some error message")
	})

	t.Run("unknown response type", func(t *testing.T) {
		raw := "not a map"
		var res User
		err := AssertResponse(raw, &res)
		require.Error(t, err)
		assert.EqualError(t, err, UnKnownResponseType)
	})
}
