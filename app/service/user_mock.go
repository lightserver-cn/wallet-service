package service

import (
	"server/app/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// MockUserRepo is a mock implementation of the repository.UserInter interface
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx *gin.Context, user *model.User) (*model.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) UpdateUser(ctx *gin.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByID(ctx *gin.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByUsername(ctx *gin.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(ctx *gin.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*model.User), args.Error(1)
}
