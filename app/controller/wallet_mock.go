package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

// MockWalletInter is a mock implementation of WalletInter
type MockWalletInter struct {
	mock.Mock
}

func (m *MockWalletInter) Deposit(ctx *gin.Context, uid int64, amount decimal.Decimal) error {
	args := m.Called(ctx, uid, amount)
	return args.Error(0)
}

func (m *MockWalletInter) Withdraw(ctx *gin.Context, uid int64, amount decimal.Decimal) error {
	args := m.Called(ctx, uid, amount)
	return args.Error(0)
}

func (m *MockWalletInter) Transfer(ctx *gin.Context, fromUID, toUID int64, amount decimal.Decimal) error {
	args := m.Called(ctx, fromUID, toUID, amount)
	return args.Error(0)
}

func (m *MockWalletInter) Balance(ctx *gin.Context, uid int64) (decimal.Decimal, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}
