package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestResponse(t *testing.T) {
	defer goleak.VerifyNone(t)

	gin.SetMode(gin.TestMode)

	t.Run("Header", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		resp := NewResponse(ctx)
		resp.Header()

		//nolint:testifylint // This is a constant string, direct comparison is sufficient, JSONEq is not needed
		assert.Equal(t, ContentTypeJSON, w.Header().Get(ContentTypeKey))
	})

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		resp := NewResponse(ctx)
		resp.Success()

		assert.Equal(t, http.StatusOK, w.Code)

		var result ResponseEntity
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, 0, result.ErrCode)
		assert.Equal(t, "success", result.ErrMsg)
		assert.Equal(t, map[string]any{}, result.Data)
	})

	t.Run("SuccessData", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		data := map[string]any{"key": "value"}

		resp := NewResponse(ctx)
		resp.SuccessData(data)

		assert.Equal(t, http.StatusOK, w.Code)

		var result ResponseEntity
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, 0, result.ErrCode)
		assert.Equal(t, "success", result.ErrMsg)
		assert.Equal(t, data, result.Data)
	})

	t.Run("SuccessMsg", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		msg := "custom success message"

		resp := NewResponse(ctx)
		resp.SuccessMsg(msg)

		assert.Equal(t, http.StatusOK, w.Code)

		var result ResponseEntity
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, 0, result.ErrCode)
		assert.Equal(t, msg, result.ErrMsg)
		assert.Equal(t, map[string]any{}, result.Data)
	})

	t.Run("SuccessDataMsg", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		data := map[string]any{"key": "value"}
		msg := "custom success message with data"

		resp := NewResponse(ctx)
		resp.SuccessDataMsg(data, msg)

		assert.Equal(t, http.StatusOK, w.Code)

		var result ResponseEntity
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, 0, result.ErrCode)
		assert.Equal(t, msg, result.ErrMsg)
		assert.Equal(t, data, result.Data)
	})

	t.Run("ErrorCodeMsg", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		errCode := 400
		msg := "bad request"
		resp := NewResponse(ctx)
		resp.ErrorCodeMsg(errCode, msg)

		assert.Equal(t, http.StatusOK, w.Code)

		var result ResponseEntity
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, errCode, result.ErrCode)
		assert.Equal(t, msg, result.ErrMsg)
		assert.Nil(t, result.Data)
	})

	t.Run("ErrorValidator", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		err := gin.Error{
			Type: gin.ErrorTypeBind,
			Err:  fmt.Errorf("validation error"),
		}
		resp := NewResponse(ctx)
		resp.ErrorValidator(err.Err)

		assert.Equal(t, http.StatusOK, w.Code)

		var result ResponseEntity
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, ErrCodeValidateErr, result.ErrCode)
		assert.Equal(t, "validation error", result.ErrMsg)
		assert.Nil(t, result.Data)
	})
}
