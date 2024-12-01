package service

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"server/app/model"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestWalletServ_NewWallet(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("TestNewWallet", func(t *testing.T) {
		// Create a mock instance
		repo := new(MockWalletRepo)

		inter := NewWallet(repo)
		assert.NotNil(t, inter)

		serv, ok := inter.(*WalletServ)
		assert.True(t, ok)
		assert.Equal(t, repo, serv.repo)
	})

	t.Run("TestNewWallet_NilRepo", func(t *testing.T) {
		inter := NewWallet(nil)
		expectedInter := &WalletServ{repo: nil}
		assert.Equal(t, expectedInter, inter)
	})
}

func TestWalletServ_Deposit(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Create a Gin context
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	mockRepo := new(MockWalletRepo)
	walletServ := NewWallet(mockRepo)

	uid := int64(1)
	amount := decimal.NewFromInt(100)

	// Mock the Balance method
	mockRepo.On("Balance", ctx, uid).Return(decimal.Zero, nil)

	mockRepo.On("Deposit", ctx, uid, amount).Return(nil)

	err := walletServ.Deposit(ctx, uid, amount)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestWalletServ_Withdraw(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Create a Gin context
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	mockRepo := new(MockWalletRepo)
	walletServ := NewWallet(mockRepo)

	uid := int64(1)
	amount := decimal.NewFromInt(100)

	// Mock the Balance method
	mockRepo.On("Balance", ctx, uid).Return(decimal.NewFromInt(500), nil)

	mockRepo.On("Withdraw", ctx, uid, amount).Return(nil)

	err := walletServ.Withdraw(ctx, uid, amount)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestWalletServ_Transfer(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Create a Gin context
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	mockRepo := new(MockWalletRepo)
	walletServ := NewWallet(mockRepo)

	fromUID := int64(1)
	toUID := int64(2)
	amount := decimal.NewFromInt(100)

	// Mock the Balance method for sender and receiver
	mockRepo.On("Balance", ctx, fromUID).Return(decimal.NewFromInt(500), nil)
	mockRepo.On("Balance", ctx, toUID).Return(decimal.NewFromInt(200), nil)

	// Mock the Transfer method
	mockRepo.On("Transfer", ctx, fromUID, toUID, amount).Return(nil)

	err := walletServ.Transfer(ctx, fromUID, toUID, amount)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestWalletServ_Balance(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Create a Gin context
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	mockRepo := new(MockWalletRepo)
	walletServ := NewWallet(mockRepo)

	uid := int64(1)

	// Mock the Balance method
	mockRepo.On("Balance", ctx, uid).Return(decimal.NewFromInt(500), nil)

	res, err := walletServ.Balance(ctx, uid)
	require.NoError(t, err)
	assert.Equal(t, decimal.NewFromInt(500), res)

	mockRepo.AssertExpectations(t)
}

func TestWalletServ_Deposit_Overflow(t *testing.T) {
	defer goleak.VerifyNone(t)

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	mockRepo := new(MockWalletRepo)
	walletServ := NewWallet(mockRepo)

	uid := int64(1)
	// Use a large amount that, when added to the near-max balance, will exceed the limit
	amount := decimal.NewFromInt(2) // Adjust this based on MaxBalance and the starting balance
	expectedErr := fmt.Errorf("deposit would exceed the maximum allowed balance of %s", decimal.NewFromInt(model.MaxBalance).String())

	// Mock the Balance method to return a value close to the maximum balance
	mockRepo.On("Balance", ctx, uid).Return(decimal.NewFromInt(model.MaxBalance-1), nil)

	// The Deposit method should not be called because the check in Deposit function should prevent it
	err := walletServ.Deposit(ctx, uid, amount)
	require.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

func TestWalletServ_Withdraw_Overflow(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Create a Gin context
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	mockRepo := new(MockWalletRepo)
	walletServ := NewWallet(mockRepo)

	uid := int64(1)
	amount := decimal.NewFromInt(1000000000000000000) // Large amount to cause overflow
	expectedErr := fmt.Errorf("insufficient balance for withdrawal")

	// Mock the Balance method to return a value close to the maximum balance
	mockRepo.On("Balance", ctx, uid).Return(decimal.NewFromInt(model.MaxBalance-1), nil)

	err := walletServ.Withdraw(ctx, uid, amount)
	require.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

func TestWalletServ_Transfer_OverflowReceiver(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Create a Gin context
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	mockRepo := new(MockWalletRepo)
	walletServ := NewWallet(mockRepo)

	fromUID := int64(1)
	toUID := int64(2)
	amount := decimal.NewFromInt(1000000000000000000) // Large amount to cause overflow
	expectedErr := fmt.Errorf("insufficient balance for transfer")

	// Mock the Balance method for sender and receiver
	mockRepo.On("Balance", ctx, fromUID).Return(decimal.NewFromInt(500), nil)

	err := walletServ.Transfer(ctx, fromUID, toUID, amount)
	require.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

func TestWalletServ_Balance_Error(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Create a Gin context
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	mockRepo := new(MockWalletRepo)
	walletServ := NewWallet(mockRepo)

	uid := int64(1)

	// Mock the Balance method to return an error
	mockRepo.On("Balance", ctx, uid).Return(decimal.Zero, fmt.Errorf("balance query error"))

	_, err := walletServ.Balance(ctx, uid)
	require.Error(t, err)
	assert.Equal(t, "balance query error", err.Error())

	mockRepo.AssertExpectations(t)
}
