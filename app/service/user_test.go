package service

import (
	"net/http/httptest"
	"testing"

	"server/app/model"
	"server/app/request"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestUserServ_NewUser(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("TestNewUser", func(t *testing.T) {
		// Create a mock instance
		repo := new(MockUserRepo)
		repoWallet := new(MockWalletRepo)

		inter := NewUser(repo, repoWallet)
		assert.NotNil(t, inter)

		serv, ok := inter.(*UserServ)
		assert.True(t, ok)
		assert.Equal(t, repo, serv.repo)
	})

	t.Run("TestNewUser_NilRepo", func(t *testing.T) {
		inter := NewUser(nil, nil)
		expectedInter := &UserServ{repo: nil, repoWallet: nil}
		assert.Equal(t, expectedInter, inter)
	})
}

func TestUserServ_RegisterUser(t *testing.T) {
	defer goleak.VerifyNone(t)

	tests := []struct {
		name          string
		req           *request.ReqRegisterUser
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "Valid Registration",
			req: &request.ReqRegisterUser{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "password123",
			},
			expectedUser: &model.User{
				ID:           1,
				Username:     "testuser",
				Email:        "testuser@example.com",
				PasswordHash: []byte("$2a$10$vI8aWBnW3fID.ZQ4/zo1GqOJQGEp9Z5X6VgKbP0F.m7fE2cUZIySsY5X3r5ZuHkK"),
				Status:       model.UserStatusInvalid,
			},
			expectedError: nil,
		},
		// Add more test cases for different scenarios
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

			mockRepo := new(MockUserRepo)
			mockWalletRepo := new(MockWalletRepo)

			userServ := NewUser(mockRepo, mockWalletRepo)

			mockRepo.On("CreateUser", ctx, mock.Anything).Return(tt.expectedUser, tt.expectedError)
			mockWalletRepo.On("CreateWallet", ctx, mock.Anything).Return(&model.Wallet{UID: tt.expectedUser.ID, Balance: decimal.NewFromFloat(0)}, nil)

			user, err := userServ.RegisterUser(ctx, tt.req)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedUser, user)

			mockRepo.AssertExpectations(t)
			mockWalletRepo.AssertExpectations(t)
		})
	}
}

func TestUserServ_UpdateUser(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Create a Gin context
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	mockRepo := new(MockUserRepo)

	userServ := NewUser(mockRepo, nil)

	user := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "testuser@example.com",
	}

	mockRepo.On("UpdateUser", ctx, user).Return(nil)

	err := userServ.UpdateUser(ctx, user)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserServ_GetUserByID(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Create a Gin context
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	mockRepo := new(MockUserRepo)

	userServ := NewUser(mockRepo, nil)

	expectedUser := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "testuser@example.com",
	}

	mockRepo.On("GetUserByID", ctx, int64(1)).Return(expectedUser, nil)

	user, err := userServ.GetUserByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	mockRepo.AssertExpectations(t)
}

func TestUserServ_GetUserByUsername(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Create a Gin context
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	mockRepo := new(MockUserRepo)
	userServ := NewUser(mockRepo, nil)

	expectedUser := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "testuser@example.com",
	}

	mockRepo.On("GetUserByUsername", ctx, "testuser").Return(expectedUser, nil)

	user, err := userServ.GetUserByUsername(ctx, "testuser")
	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	mockRepo.AssertExpectations(t)
}

func TestUserServ_GetUserByEmail(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Create a Gin context
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	mockRepo := new(MockUserRepo)
	userServ := NewUser(mockRepo, nil)

	expectedUser := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "testuser@example.com",
	}

	mockRepo.On("GetUserByEmail", ctx, "testuser@example.com").Return(expectedUser, nil)

	user, err := userServ.GetUserByEmail(ctx, "testuser@example.com")
	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_PasswordEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	userServ := &UserServ{}

	req := &request.ReqRegisterUser{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "",
	}

	_, err := userServ.RegisterUser(ctx, req)

	require.Error(t, err)
	assert.Equal(t, "password cannot be empty", err.Error())
}

func TestRegisterUser_EmailEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	userServ := &UserServ{}

	req := &request.ReqRegisterUser{
		Username: "testuser",
		Email:    "",
		Password: "password123",
	}

	_, err := userServ.RegisterUser(ctx, req)

	require.Error(t, err)
	assert.Equal(t, "email cannot be empty", err.Error())
}
