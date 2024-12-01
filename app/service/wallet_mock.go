package service

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"server/app/model"
)

// MockWalletRepo is a mock implementation of the repository.WalletInter interface
type MockWalletRepo struct {
	mock.Mock
}

func (m *MockWalletRepo) CreateWallet(ctx *gin.Context, wallet *model.Wallet) (*model.Wallet, error) {
	args := m.Called(ctx, wallet)
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWalletRepo) GetWalletByUID(ctx *gin.Context, uid int64) (*model.Wallet, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWalletRepo) Deposit(ctx *gin.Context, uid int64, amount decimal.Decimal) error {
	args := m.Called(ctx, uid, amount)
	return args.Error(0)
}

func (m *MockWalletRepo) Withdraw(ctx *gin.Context, uid int64, amount decimal.Decimal) error {
	args := m.Called(ctx, uid, amount)
	return args.Error(0)
}

func (m *MockWalletRepo) Transfer(ctx *gin.Context, fromUID, toUID int64, amount decimal.Decimal) error {
	args := m.Called(ctx, fromUID, toUID, amount)
	return args.Error(0)
}

func (m *MockWalletRepo) Balance(ctx *gin.Context, uid int64) (decimal.Decimal, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}
