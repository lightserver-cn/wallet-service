package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"server/app/model"
	"server/app/request"
	"server/pkg/consts"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

// Test cases for UserCtrl.RegisterUser
func TestUserCtrl_RegisterUser(t *testing.T) {
	defer goleak.VerifyNone(t)

	gin.SetMode(gin.TestMode)
	mockService := new(MockUserInter)
	userCtrl := NewUser(mockService)

	tests := []struct {
		name                      string
		req                       *request.ReqRegisterUser
		mockGetUserByUsername     *model.User
		mockGetUserByEmail        *model.User
		mockRegisterUser          *model.User
		mockGetUserByUsernameSkip bool
		mockGetUserByEmailSkip    bool
		mockRegisterUserSkip      bool
		mockGetUserByUsernameErr  error
		mockGetUserByEmailErr     error
		mockRegisterErr           error
		expectedStatus            int
		expectedError             string
	}{
		{
			name: "Valid user registration",
			req: &request.ReqRegisterUser{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			mockRegisterUser: &model.User{
				ID:        2,
				Username:  "newuser",
				Email:     "newuser@example.com",
				Status:    1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Username already exists",
			req: &request.ReqRegisterUser{
				Username: "existinguser",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			mockGetUserByUsername: &model.User{
				ID:        1,
				Username:  "existinguser",
				Email:     "existinguser@example.com",
				Status:    1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockGetUserByEmailSkip: true,
			mockRegisterUserSkip:   true,
			expectedStatus:         http.StatusConflict,
			expectedError:          consts.ErrUsernameAlreadyExists,
		},
		{
			name: "Email already exists",
			req: &request.ReqRegisterUser{
				Username: "newuser",
				Email:    "existinguser@example.com",
				Password: "password123",
			},
			mockGetUserByEmail: &model.User{
				ID:        1,
				Username:  "existinguser",
				Email:     "existinguser@example.com",
				Status:    1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockRegisterUserSkip: true,
			expectedStatus:       http.StatusConflict,
			expectedError:        consts.ErrEmailAlreadyExists,
		},
		{
			name: "Internal server error on GetUserByUsername",
			req: &request.ReqRegisterUser{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			mockGetUserByEmailSkip:   true,
			mockRegisterUserSkip:     true,
			mockGetUserByUsernameErr: errors.New(consts.ErrInternalServer),
			expectedStatus:           http.StatusInternalServerError,
			expectedError:            consts.ErrInternalServer,
		},
		{
			name: "Internal server error on GetUserByEmail",
			req: &request.ReqRegisterUser{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			mockRegisterUserSkip:  true,
			mockGetUserByEmailErr: errors.New(consts.ErrInternalServer),
			expectedStatus:        http.StatusInternalServerError,
			expectedError:         consts.ErrInternalServer,
		},
		{
			name: "Internal server error on RegisterUser",
			req: &request.ReqRegisterUser{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			mockRegisterErr: errors.New(consts.ErrInternalServer),
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   consts.ErrInternalServer,
		},
		{
			name: "Validation username failed",
			req: &request.ReqRegisterUser{
				Username: "", // Invalid username
				Email:    "newuser@example.com",
				Password: "password123",
			},
			mockGetUserByUsernameSkip: true,
			mockGetUserByEmailSkip:    true,
			mockRegisterUserSkip:      true,
			expectedStatus:            http.StatusBadRequest,
			expectedError:             consts.ErrUsernameRequired,
		},
		{
			name: "Validation email failed",
			req: &request.ReqRegisterUser{
				Username: "newuser", // Invalid username
				Email:    "",
				Password: "password123",
			},
			mockGetUserByUsernameSkip: true,
			mockGetUserByEmailSkip:    true,
			mockRegisterUserSkip:      true,
			expectedStatus:            http.StatusBadRequest,
			expectedError:             consts.ErrEmailRequired,
		},
		{
			name: "Validation password failed",
			req: &request.ReqRegisterUser{
				Username: "newuser", // Invalid username
				Email:    "newuser@example.com",
				Password: "",
			},
			mockGetUserByUsernameSkip: true,
			mockGetUserByEmailSkip:    true,
			mockRegisterUserSkip:      true,
			expectedStatus:            http.StatusBadRequest,
			expectedError:             consts.ErrPasswordRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			if !tt.mockGetUserByUsernameSkip {
				mockService.On("GetUserByUsername", ctx, tt.req.Username).Return(tt.mockGetUserByUsername, tt.mockGetUserByUsernameErr)
			}

			if !tt.mockGetUserByEmailSkip {
				mockService.On("GetUserByEmail", ctx, tt.req.Email).Return(tt.mockGetUserByEmail, tt.mockGetUserByEmailErr)
			}

			if !tt.mockRegisterUserSkip {
				mockService.On("RegisterUser", ctx, tt.req).Return(tt.mockRegisterUser, tt.mockRegisterErr)
			}

			reqBody, err := json.Marshal(tt.req)
			require.NoError(t, err)
			ctx.Request, err = http.NewRequest("POST", "", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			ctx.Request.Header.Set("Content-Type", "application/json")

			userCtrl.RegisterUser(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// Test cases for UserCtrl.GetUserByUID
func TestUserCtrlGetUserByUID(t *testing.T) {
	defer goleak.VerifyNone(t)

	gin.SetMode(gin.TestMode)
	mockService := new(MockUserInter)
	userCtrl := NewUser(mockService)

	tests := []struct {
		name            string
		uid             int64
		mockGetUser     *model.User
		mockGetUserSkip bool
		mockGetUserErr  error
		expectedStatus  int
		expectedError   string
	}{
		{
			name: "Valid user retrieval",
			uid:  1,
			mockGetUser: &model.User{
				ID:        1,
				Username:  "testuser",
				Email:     "testuser@example.com",
				Status:    1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User not found",
			uid:            999,
			mockGetUserErr: sql.ErrNoRows,
			expectedStatus: http.StatusNotFound,
			expectedError:  consts.ErrUserNotFound,
		},
		{
			name:            "Invalid UID",
			uid:             0,
			mockGetUserSkip: true,
			mockGetUserErr:  errors.New(consts.ErrInternalServer),
			expectedStatus:  http.StatusBadRequest,
			expectedError:   consts.ErrInvalidUID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = gin.Params{
				{Key: "uid", Value: strconv.FormatInt(tt.uid, 10)},
			}

			if !tt.mockGetUserSkip {
				mockService.On("GetUserByID", ctx, tt.uid).Return(tt.mockGetUser, tt.mockGetUserErr)
			}

			userCtrl.GetUserByUID(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}
