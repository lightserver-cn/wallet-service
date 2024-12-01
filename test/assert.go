package test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
)

const UnKnownResponseType = "unknown response type"

type SuccessResponse struct {
	Message string `json:"message"`
}

// AssertResponseSuccess struct is used to parse an success response.
func AssertResponseSuccess(t *testing.T, expected string, response *httpexpect.Value, msgAndArgs ...any) {
	assert.Equal(t, expected, response.Path("$.message").String().Raw(), msgAndArgs)
}

// ErrorResponse struct is used to parse an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

func AssertResponseError(t *testing.T, expected string, response *httpexpect.Value, msgAndArgs ...any) {
	assert.Equal(t, expected, response.Path("$.error").String().Raw(), msgAndArgs)
}

func AssertResponse(raw, res any) error {
	// Examine the type of response and handle it accordingly.
	switch v := raw.(type) {
	case map[string]any:
		// Check if there is an "error" field in the response.
		if errMsg, ok := v["error"].(string); ok {
			return errors.New(errMsg)
		} else {
			userBytes, _ := json.Marshal(v) // Convert the map to JSON bytes.
			return json.Unmarshal(userBytes, res)
		}
	default:
		return errors.New(UnKnownResponseType)
	}
}
