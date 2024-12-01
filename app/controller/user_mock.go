package controller

import (
	"server/app/model"
	"server/app/request"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// MockUserInter is a mock implementation of UserInter
type MockUserInter struct {
	mock.Mock
}

func (m *MockUserInter) RegisterUser(ctx *gin.Context, req *request.ReqRegisterUser) (*model.User, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserInter) UpdateUser(ctx *gin.Context, mod *model.User) error {
	args := m.Called(ctx, mod)
	return args.Error(0)
}

func (m *MockUserInter) GetUserByID(ctx *gin.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserInter) GetUserByUsername(ctx *gin.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserInter) GetUserByEmail(ctx *gin.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*model.User), args.Error(1)
}
