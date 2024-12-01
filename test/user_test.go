package test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"server/app/model"
	"server/pkg/consts"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestPostUsers(t *testing.T) {
	defer goleak.VerifyNone(
		t,
		goleak.IgnoreTopFunction("net/http.(*Server).Serve"),
		goleak.IgnoreTopFunction("net/http/httptest.(*Server).goServe.func1"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
		goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
		goleak.IgnoreTopFunction("internal/poll.(*pollDesc).wait"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Accept"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Read"),
		goleak.IgnoreTopFunction("time.Sleep"),
		goleak.IgnoreTopFunction("time.AfterFunc"),
		goleak.IgnoreTopFunction("time.Ticker"),
		goleak.IgnoreTopFunction("runtime.gopark"),
		goleak.IgnoreTopFunction("runtime.forcegchelper"),
		goleak.IgnoreTopFunction("runtime.bgsweep"),
		goleak.IgnoreTopFunction("runtime.bgscavenge"),
	)

	m := NewMockTest().start(t)
	defer m.Teardown()

	t.Run("register-user-success", func(t *testing.T) {
		req := map[string]any{
			"username": "TestPostUsers",
			"email":    "TestPostUsers@gmail.com",
			"password": "TestPostUsers",
		}
		res := m.Expect.POST("/api/users").WithJSON(req).Expect().Status(http.StatusCreated).JSON()

		resp := &model.User{}
		if err := AssertResponse(res.Raw(), &resp); err != nil {
			t.Error(err)
		}

		assert.Equal(t, req["username"], resp.Username, "username mismatch")
		assert.Equal(t, req["email"], resp.Email, "email mismatch")
	})

	t.Run("register-user-empty-username", func(t *testing.T) {
		req := map[string]any{
			"username": "",
			"email":    "TestPostUsers@gmail.com",
			"password": "TestPostUsers",
		}

		res := m.Expect.POST("/api/users").WithJSON(req).Expect().Status(http.StatusBadRequest).JSON()

		AssertResponseError(t, consts.ErrUsernameRequired, res, "error mismatch")
	})

	t.Run("register-user-empty-email", func(t *testing.T) {
		req := map[string]any{
			"username": "TestPostUsers",
			"email":    "",
			"password": "TestPostUsers",
		}

		res := m.Expect.POST("/api/users").WithJSON(req).Expect().Status(http.StatusBadRequest).JSON()

		AssertResponseError(t, consts.ErrEmailRequired, res, "error mismatch")
	})

	t.Run("register-user-empty-password", func(t *testing.T) {
		req := map[string]any{
			"username": "TestPostUsers",
			"email":    "TestPostUsers@gmail.com",
			"password": "",
		}

		res := m.Expect.POST("/api/users").WithJSON(req).Expect().Status(http.StatusBadRequest).JSON()

		AssertResponseError(t, consts.ErrPasswordRequired, res, "error mismatch")
	})

	t.Run("register-user-username-repeated", func(t *testing.T) {
		req := map[string]any{
			"username": "Test1PostUsers",
			"email":    "Test1PostUsers@gmail.com",
			"password": "Test1PostUsers",
		}

		res := m.Expect.POST("/api/users").WithJSON(req).Expect().Status(http.StatusCreated).JSON()

		resp := &model.User{}
		if err := AssertResponse(res.Raw(), &resp); err != nil {
			t.Error(err)
		}

		assert.Equal(t, req["username"], resp.Username, "username mismatch")
		assert.Equal(t, req["email"], resp.Email, "email mismatch")

		req["email"] = "Test2PostUsers@gmail.com"
		res = m.Expect.POST("/api/users").WithJSON(req).Expect().Status(http.StatusConflict).JSON()

		AssertResponseError(t, consts.ErrUsernameAlreadyExists, res, "error mismatch")
	})

	t.Run("register-user-email-repeated", func(t *testing.T) {
		req := map[string]any{
			"username": "Test2PostUsers",
			"email":    "Test2PostObjects@gmail.com",
			"password": "Test2PostUsers",
		}

		res := m.Expect.POST("/api/users").WithJSON(req).Expect().Status(http.StatusCreated).JSON()

		resp := &model.User{}
		if err := AssertResponse(res.Raw(), &resp); err != nil {
			t.Error(err)
		}

		assert.Equal(t, req["username"], resp.Username, "username mismatch")
		assert.Equal(t, req["email"], resp.Email, "email mismatch")

		req["username"] = "Test3PostUsers"
		res = m.Expect.POST("/api/users").WithJSON(req).Expect().Status(http.StatusConflict).JSON()

		AssertResponseError(t, consts.ErrEmailAlreadyExists, res, "error mismatch")
	})
}

func TestGetUsersByUID(t *testing.T) {
	defer goleak.VerifyNone(
		t,
		goleak.IgnoreTopFunction("net/http.(*Server).Serve"),
		goleak.IgnoreTopFunction("net/http/httptest.(*Server).goServe.func1"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
		goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
		goleak.IgnoreTopFunction("internal/poll.(*pollDesc).wait"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Accept"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Read"),
		goleak.IgnoreTopFunction("time.Sleep"),
		goleak.IgnoreTopFunction("time.AfterFunc"),
		goleak.IgnoreTopFunction("time.Ticker"),
		goleak.IgnoreTopFunction("runtime.gopark"),
		goleak.IgnoreTopFunction("runtime.forcegchelper"),
		goleak.IgnoreTopFunction("runtime.bgsweep"),
		goleak.IgnoreTopFunction("runtime.bgscavenge"),
	)

	m := NewMockTest().start(t)
	defer m.Teardown()

	t.Run("get-user-by-uid-success", func(t *testing.T) {
		req := map[string]string{
			"username": "TestGetUsersByUID",
			"email":    "TestGetUsersByUID@gmail.com",
			"password": "TestGetUsersByUI",
		}

		resRegister := m.Expect.POST("/api/users").WithJSON(req).Expect().Status(http.StatusCreated).JSON()

		respRegister := &model.User{}
		if err := AssertResponse(resRegister.Raw(), &respRegister); err != nil {
			t.Error(err)
		}

		assert.Equal(t, req["username"], respRegister.Username, "username mismatch")
		assert.Equal(t, req["email"], respRegister.Email, "email mismatch")

		uid := respRegister.ID

		resGetByUID := m.Expect.GET(fmt.Sprintf("/api/users/%d", uid)).
			Expect().Status(http.StatusOK).JSON()

		respGetByUID := &model.User{}
		if err := AssertResponse(resGetByUID.Raw(), &respGetByUID); err != nil {
			t.Error(err)
		}

		assert.Equal(t, req["username"], respGetByUID.Username, "username mismatch")
		assert.Equal(t, req["email"], respGetByUID.Email, "email mismatch")
	})

	t.Run("get-user-by-uid-not-exist", func(t *testing.T) {
		var uidNotExist int64 = 9999
		resGetByUID := m.Expect.GET("/api/users/" + strconv.FormatInt(uidNotExist, 10)).Expect().Status(http.StatusNotFound).JSON()

		AssertResponseError(t, consts.ErrUserNotFound, resGetByUID, "error mismatch")
	})

	t.Run("get-user-by-uid-validation-failed", func(t *testing.T) {
		var uidErr = "aaa"
		resGetByUI := m.Expect.GET("/api/users/" + uidErr).Expect().Status(http.StatusBadRequest).JSON()

		AssertResponseError(t, consts.ErrValidationFailed, resGetByUI, "error mismatch")
	})
}
