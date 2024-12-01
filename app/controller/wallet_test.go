package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"server/app/model"
	"server/app/request"
	"server/pkg/consts"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

// Test cases for WalletCtrl.Deposit
func TestWalletCtrl_Deposit(t *testing.T) {
	defer goleak.VerifyNone(t)

	gin.SetMode(gin.TestMode)
	mockService := new(MockWalletInter)
	walletCtrl := NewWallet(mockService, nil) // Assuming NewWallet only needs WalletInter for transactions

	tests := []struct {
		name            string
		uid             int64
		amount          decimal.Decimal
		mockDepositSkip bool
		mockDepositErr  error
		expectedStatus  int
		expectedError   string
	}{
		{
			name:           "Valid deposit",
			uid:            1,
			amount:         decimal.NewFromInt(100),
			mockDepositErr: nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:            "Invalid UID",
			uid:             -1,
			amount:          decimal.NewFromInt(100),
			mockDepositSkip: true,
			expectedStatus:  http.StatusBadRequest,
			expectedError:   consts.ErrInvalidUID,
		},
		{
			name:            "Invalid Amount",
			uid:             1,
			amount:          decimal.NewFromInt(0),
			mockDepositSkip: true,
			expectedStatus:  http.StatusBadRequest,
			expectedError:   consts.ErrInvalidAmount,
		},
		{
			name:           consts.ErrInternalServer,
			uid:            1,
			amount:         decimal.NewFromInt(100),
			mockDepositErr: errors.New(consts.ErrInternalServer),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  consts.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = gin.Params{
				{Key: "uid", Value: strconv.FormatInt(tt.uid, 10)},
			}

			reqBody, err := json.Marshal(&request.ReqAmount{
				Amount: tt.amount,
			})
			require.NoError(t, err)
			ctx.Request, err = http.NewRequest("POST", "", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			ctx.Request.Header.Set("Content-Type", "application/json")

			if !tt.mockDepositSkip {
				mockService.On("Deposit", ctx, tt.uid, tt.amount).Return(tt.mockDepositErr)
			}

			walletCtrl.Deposit(ctx)

			assert.Equal(t, tt.expectedStatus, ctx.Writer.Status())

			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// Test cases for WalletCtrl.Withdraw
func TestWalletCtrl_Withdraw(t *testing.T) {
	defer goleak.VerifyNone(t)

	gin.SetMode(gin.TestMode)
	mockService := new(MockWalletInter)
	walletCtrl := NewWallet(mockService, nil) // Assuming NewWallet only needs WalletInter for transactions

	tests := []struct {
		name             string
		uid              int64
		amount           decimal.Decimal
		mockWithdrawSkip bool
		mockWithdrawErr  error
		expectedStatus   int
		expectedError    string
	}{
		{
			name:            "Valid withdraw",
			uid:             1,
			amount:          decimal.NewFromInt(100),
			mockWithdrawErr: nil,
			expectedStatus:  http.StatusOK,
		},
		{
			name:             "Invalid UID",
			uid:              -1,
			amount:           decimal.NewFromInt(100),
			mockWithdrawSkip: true,
			expectedStatus:   http.StatusBadRequest,
			expectedError:    consts.ErrInvalidUID,
		},
		{
			name:             "Invalid Amount",
			uid:              1,
			amount:           decimal.NewFromInt(0),
			mockWithdrawSkip: true,
			expectedStatus:   http.StatusBadRequest,
			expectedError:    consts.ErrInvalidAmount,
		},
		{
			name:            "Insufficient balance",
			uid:             1,
			amount:          decimal.NewFromInt(1000),
			mockWithdrawErr: errors.New("insufficient balance"),
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   "insufficient balance",
		},
		{
			name:            consts.ErrInternalServer,
			uid:             1,
			amount:          decimal.NewFromInt(100),
			mockWithdrawErr: errors.New(consts.ErrInternalServer),
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   consts.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = gin.Params{
				{Key: "uid", Value: strconv.FormatInt(tt.uid, 10)},
			}

			reqBody, err := json.Marshal(&request.ReqAmount{
				Amount: tt.amount,
			})
			require.NoError(t, err)
			ctx.Request, err = http.NewRequest("POST", "", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			ctx.Request.Header.Set("Content-Type", "application/json")

			if !tt.mockWithdrawSkip {
				mockService.On("Withdraw", ctx, tt.uid, tt.amount).Return(tt.mockWithdrawErr)
			}

			walletCtrl.Withdraw(ctx)

			assert.Equal(t, tt.expectedStatus, ctx.Writer.Status())

			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// Test cases for WalletCtrl.Transfer
func TestWalletCtrl_Transfer(t *testing.T) {
	defer goleak.VerifyNone(t)

	gin.SetMode(gin.TestMode)
	mockService := new(MockWalletInter)
	walletCtrl := NewWallet(mockService, nil) // Assuming NewWallet only needs WalletInter for transactions

	tests := []struct {
		name             string
		uid              int64
		toUID            int64
		amount           decimal.Decimal
		mockTransfer     error
		mockTransferSkip bool
		expectedStatus   int
		expectedError    string
	}{
		{
			name:           "Valid transfer",
			uid:            1,
			toUID:          2,
			amount:         decimal.NewFromInt(100),
			mockTransfer:   nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:             "Invalid UID",
			uid:              -1,
			toUID:            2,
			amount:           decimal.NewFromInt(100),
			mockTransferSkip: true,
			expectedStatus:   http.StatusBadRequest,
			expectedError:    consts.ErrInvalidUID,
		},
		{
			name:             "Invalid ToUID",
			uid:              1,
			toUID:            -1,
			amount:           decimal.NewFromInt(100),
			mockTransferSkip: true,
			expectedStatus:   http.StatusBadRequest,
			expectedError:    consts.ErrInvalidUID,
		},
		{
			name:             "Invalid Amount",
			uid:              1,
			toUID:            2,
			amount:           decimal.NewFromInt(0),
			mockTransferSkip: true,
			expectedStatus:   http.StatusBadRequest,
			expectedError:    consts.ErrInvalidAmount,
		},
		{
			name:           consts.ErrInternalServer,
			uid:            1,
			toUID:          2,
			amount:         decimal.NewFromInt(100),
			mockTransfer:   errors.New(consts.ErrInternalServer),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  consts.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = gin.Params{
				{Key: "uid", Value: strconv.FormatInt(tt.uid, 10)},
			}

			reqBody, err := json.Marshal(&request.ReqTransfer{
				ToUID:  tt.toUID,
				Amount: tt.amount,
			})
			require.NoError(t, err)
			ctx.Request, err = http.NewRequest("POST", "", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			ctx.Request.Header.Set("Content-Type", "application/json")

			if !tt.mockTransferSkip {
				mockService.On("Transfer", ctx, tt.uid, tt.toUID, tt.amount).Return(tt.mockTransfer)
			}

			walletCtrl.Transfer(ctx)

			assert.Equal(t, tt.expectedStatus, ctx.Writer.Status())

			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// Test cases for WalletCtrl.Balance
func TestWalletCtrl_Balance(t *testing.T) {
	defer goleak.VerifyNone(t)

	gin.SetMode(gin.TestMode)
	mockService := new(MockWalletInter)
	walletCtrl := NewWallet(mockService, nil) // Assuming NewWallet only needs WalletInter for transactions

	tests := []struct {
		name            string
		uid             int64
		mockBalance     decimal.Decimal
		mockBalanceSkip bool
		mockBalanceErr  error
		expectedStatus  int
		expectedError   string
	}{
		{
			name:           "Valid balance retrieval",
			uid:            1,
			mockBalance:    decimal.NewFromInt(100),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User not found",
			uid:            2,
			mockBalanceErr: sql.ErrNoRows,
			expectedStatus: http.StatusNotFound,
			expectedError:  consts.ErrUserNotFound,
		},
		{
			name:           consts.ErrInternalServer,
			uid:            3,
			mockBalanceErr: errors.New(consts.ErrInternalServer),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  consts.ErrInternalServer,
		},
		{
			name:            "Invalid UID",
			uid:             0,
			mockBalanceSkip: true,
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

			if !tt.mockBalanceSkip {
				mockService.On("Balance", ctx, tt.uid).
					Return(tt.mockBalance, tt.mockBalanceErr)
			}

			walletCtrl.Balance(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// Test cases for WalletCtrl.Transactions
func TestWalletCtrl_Transactions(t *testing.T) {
	defer goleak.VerifyNone(t)

	gin.SetMode(gin.TestMode)
	transactionService := new(MockTransactionInter)
	walletCtrl := NewWallet(nil, transactionService) // Assuming NewWallet only needs TransactionInter for transactions

	tests := []struct {
		name                string
		req                 *request.ReqTransactions
		mockTransaction     *request.ResTransactions
		mockTransactionErr  error
		mockTransactionSkip bool
		expectedStatus      int
		expectedError       string
	}{
		{
			name: "Valid transaction",
			req: &request.ReqTransactions{
				UID:  2,
				Type: 2,
				ReqPage: request.ReqPage{
					Page:     2,
					PageSize: 2,
				},
			},
			mockTransaction: &request.ResTransactions{
				List:    []*model.TransactionWithUsername{},
				HasMore: false,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "User not found",
			req: &request.ReqTransactions{
				UID:  2,
				Type: 2,
				ReqPage: request.ReqPage{
					Page:     2,
					PageSize: 2,
				},
			},
			mockTransactionErr: sql.ErrNoRows,
			mockTransaction: &request.ResTransactions{
				List:    []*model.TransactionWithUsername{},
				HasMore: false,
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  sql.ErrNoRows.Error(),
		},
		{
			name: "Invalid UID",
			req: &request.ReqTransactions{
				UID:  -1,
				Type: 1,
				ReqPage: request.ReqPage{
					Page:     1,
					PageSize: 10,
				},
			},
			mockTransactionSkip: true,
			expectedStatus:      http.StatusBadRequest,
			expectedError:       consts.ErrInvalidUID,
		},
		{
			name: "Invalid transaction type",
			req: &request.ReqTransactions{
				UID:  1,
				Type: 128,
				ReqPage: request.ReqPage{
					Page:     1,
					PageSize: 10,
				},
			},
			mockTransactionSkip: true,
			expectedStatus:      http.StatusBadRequest,
			expectedError:       consts.ErrInvalidTransactionType,
		},
		{
			name: "Empty result set",
			req: &request.ReqTransactions{
				UID:  1,
				Type: 1,
				ReqPage: request.ReqPage{
					Page:     1,
					PageSize: 10,
				},
			},
			mockTransaction: &request.ResTransactions{
				List:    []*model.TransactionWithUsername{},
				HasMore: false,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: consts.ErrInternalServer,
			req: &request.ReqTransactions{
				UID:  1,
				Type: 1,
				ReqPage: request.ReqPage{
					Page:     1,
					PageSize: 10,
				},
			},
			mockTransactionErr: errors.New(consts.ErrInternalServer),
			expectedStatus:     http.StatusInternalServerError,
			expectedError:      consts.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = gin.Params{
				{Key: "uid", Value: strconv.FormatInt(tt.req.UID, 10)},
			}

			if !tt.mockTransactionSkip {
				transactionService.On("GetTransactionsByUID", ctx, tt.req).
					Return(tt.mockTransaction, tt.mockTransactionErr)
			}

			var err error
			var url = fmt.Sprintf("/?page=%d&page_size=%d&type=%d", tt.req.Page, tt.req.PageSize, tt.req.Type)

			ctx.Request, err = http.NewRequest("GET", url, http.NoBody)
			require.NoError(t, err)

			walletCtrl.Transactions(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "Skip mockTransaction" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			transactionService.AssertExpectations(t)
		})
	}
}
