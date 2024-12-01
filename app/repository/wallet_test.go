package repository

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"go.uber.org/zap"

	"server/app/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestWalletRepo_NewWallet(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("TestNewWallet", func(t *testing.T) {
		db, _, errNew := sqlmock.New()
		require.NoError(t, errNew)
		defer db.Close()

		inter := NewWallet(db, zap.NewExample().Sugar())
		assert.NotNil(t, inter)

		repo, ok := inter.(*WalletRepo)
		assert.True(t, ok)
		assert.Equal(t, db, repo.db)
	})

	t.Run("TestNewWallet_NilDB", func(t *testing.T) {
		inter := NewWallet(nil, nil)
		expectedInter := &WalletRepo{db: nil}
		assert.Equal(t, expectedInter, inter)
	})
}

func TestWalletRepo_CreateWallet(t *testing.T) {
	defer goleak.VerifyNone(t)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	walletRepo := &WalletRepo{
		db:     db,
		logger: zap.NewExample().Sugar(),
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	t.Run("CreateWallet_Normal", func(t *testing.T) {
		mod := &model.Wallet{
			UID:     123,
			Balance: decimal.NewFromFloat(100.0),
		}

		expectedID := int64(1)
		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletInsert)).
			WithArgs(mod.UID, mod.Balance).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

		createdWallet, err := walletRepo.CreateWallet(ctx, mod)
		require.NoError(t, err)
		assert.Equal(t, expectedID, createdWallet.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateWallet_PrepareError", func(t *testing.T) {
		mod := &model.Wallet{
			UID:     456,
			Balance: decimal.NewFromFloat(200.0),
		}

		expectedErr := fmt.Errorf("simulated prepare error")
		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletInsert)).
			WillReturnError(expectedErr)

		createdWallet, err := walletRepo.CreateWallet(ctx, mod)
		assert.Equal(t, mod, createdWallet)
		assert.Equal(t, fmt.Errorf("failed to insert wallet: %w", expectedErr), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateWallet_QueryError", func(t *testing.T) {
		mod := &model.Wallet{
			UID:     789,
			Balance: decimal.NewFromFloat(300.0),
		}

		expectedErr := fmt.Errorf("simulated query error")
		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletInsert)).
			WithArgs(mod.UID, mod.Balance).
			WillReturnError(expectedErr)

		createdWallet, err := walletRepo.CreateWallet(ctx, mod)
		assert.Equal(t, mod, createdWallet)
		assert.Equal(t, fmt.Errorf("failed to insert wallet: %w", expectedErr), err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateWallet_EmptyResultError", func(t *testing.T) {
		mod := &model.Wallet{
			UID:     111,
			Balance: decimal.NewFromFloat(400.0),
		}

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletInsert)).
			WithArgs(mod.UID, mod.Balance).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		createdWallet, err := walletRepo.CreateWallet(ctx, mod)
		assert.Equal(t, mod, createdWallet)
		require.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWalletRepo_GetWalletByUID(t *testing.T) {
	defer goleak.VerifyNone(t)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	walletRepo := &WalletRepo{
		db:     db,
		logger: zap.NewExample().Sugar(),
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	filed := "uid"
	t.Run("GetWalletByUID_Normal", func(t *testing.T) {
		uid := int64(123)
		expectedWallet := &model.Wallet{
			ID:        1,
			UID:       uid,
			Balance:   decimal.NewFromFloat(100.5),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletByField)).
			WithArgs(filed, uid).
			WillReturnRows(sqlmock.NewRows([]string{"id", "uid", "balance", "created_at", "updated_at"}).
				AddRow(expectedWallet.ID, expectedWallet.UID, expectedWallet.Balance, expectedWallet.CreatedAt, expectedWallet.UpdatedAt))

		wallet, err := walletRepo.GetWalletByUID(ctx, uid)
		require.NoError(t, err)
		assert.Equal(t, expectedWallet, wallet)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetWalletByUID_NoRows", func(t *testing.T) {
		uid := int64(456)
		expectedErr := fmt.Errorf("sql: no rows in result set")
		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletByField)).
			WithArgs(filed, uid).
			WillReturnRows(sqlmock.NewRows([]string{"id", "uid", "balance", "created_at", "updated_at"}))

		wallet, err := walletRepo.GetWalletByUID(ctx, uid)
		assert.Equal(t, &model.Wallet{}, wallet)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetWalletByUID_PrepareError", func(t *testing.T) {
		uid := int64(789)
		expectedErr := fmt.Errorf("failed to query model by field: %w", fmt.Errorf("simulated prepare error"))
		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletByField)).
			WillReturnError(expectedErr)

		wallet, err := walletRepo.GetWalletByUID(ctx, uid)
		assert.Equal(t, &model.Wallet{}, wallet)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetWalletByUID_QueryError", func(t *testing.T) {
		uid := int64(101112)
		expectedErr := fmt.Errorf("simulated query error")
		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletByField)).
			WithArgs(filed, uid).
			WillReturnError(expectedErr)

		wallet, err := walletRepo.GetWalletByUID(ctx, uid)
		assert.Equal(t, &model.Wallet{}, wallet)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWalletRepo_queryModelByField(t *testing.T) {
	defer goleak.VerifyNone(t)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	walletRepo := &WalletRepo{
		db:     db,
		logger: zap.NewExample().Sugar(),
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	field := "uid"
	t.Run("queryModelByField_Normal", func(t *testing.T) {
		value := int64(123)
		expectedWallet := &model.Wallet{
			ID:        1,
			UID:       value,
			Balance:   decimal.NewFromFloat(100.5),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletByField)).
			WithArgs("uid", value).
			WillReturnRows(sqlmock.NewRows([]string{"id", "uid", "balance", "created_at", "updated_at"}).
				AddRow(expectedWallet.ID, expectedWallet.UID, expectedWallet.Balance, expectedWallet.CreatedAt, expectedWallet.UpdatedAt))

		wallet, err := walletRepo.queryModelByField(ctx, field, value)
		require.NoError(t, err)
		assert.Equal(t, expectedWallet, wallet)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("queryModelByField_NoRows", func(t *testing.T) {
		value := int64(456)
		expectedErr := fmt.Errorf("sql: no rows in result set")

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletByField)).
			WithArgs("uid", value).
			WillReturnRows(sqlmock.NewRows([]string{"id", "uid", "balance", "created_at", "updated_at"}))

		wallet, err := walletRepo.queryModelByField(ctx, field, value)
		assert.Equal(t, &model.Wallet{}, wallet)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("queryModelByField_PrepareError", func(t *testing.T) {
		value := int64(789)
		expectedErr := fmt.Errorf("simulated prepare error")

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletByField)).
			WillReturnError(expectedErr)

		wallet, err := walletRepo.queryModelByField(ctx, field, value)
		assert.Equal(t, &model.Wallet{}, wallet)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("queryModelByField_QueryError", func(t *testing.T) {
		value := int64(101112)
		expectedErr := fmt.Errorf("simulated query error")

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletByField)).
			WithArgs("uid", value).
			WillReturnError(expectedErr)

		wallet, err := walletRepo.queryModelByField(ctx, field, value)
		assert.Equal(t, &model.Wallet{}, wallet)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWalletRepo_Deposit(t *testing.T) {
	defer goleak.VerifyNone(t)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	walletRepo := &WalletRepo{
		db:     db,
		logger: zap.NewExample().Sugar(),
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	t.Run("TestDeposit_Success", func(t *testing.T) {
		uid := int64(123)
		amount := decimal.NewFromFloat(100.5)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletDeposit)).
			WithArgs(amount, uid, model.MaxBalance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(model.QueryInsertTransaction)).
			WithArgs(0, uid, amount, model.TransactionTypeDeposit).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := walletRepo.Deposit(ctx, uid, amount)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TestDeposit_UpdateError", func(t *testing.T) {
		uid := int64(456)
		amount := decimal.NewFromFloat(100.5)
		expectedErr := fmt.Errorf("update failed")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletDeposit)).
			WithArgs(amount, uid, model.MaxBalance).
			WillReturnError(expectedErr)
		mock.ExpectRollback()

		err := walletRepo.Deposit(ctx, uid, amount)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TestDeposit_InsertError", func(t *testing.T) {
		uid := int64(789)
		amount := decimal.NewFromFloat(100.5)
		expectedErr := fmt.Errorf("insert failed")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletDeposit)).
			WithArgs(amount, uid, model.MaxBalance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(model.QueryInsertTransaction)).
			WithArgs(0, uid, amount, model.TransactionTypeDeposit).
			WillReturnError(expectedErr)
		mock.ExpectRollback()

		err := walletRepo.Deposit(ctx, uid, amount)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWalletRepo_Withdraw(t *testing.T) {
	defer goleak.VerifyNone(t)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	walletRepo := &WalletRepo{
		db:     db,
		logger: zap.NewExample().Sugar(),
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	t.Run("TestWithdraw_Success", func(t *testing.T) {
		uid := int64(123)
		amount := decimal.NewFromFloat(100.5)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletWithdraw)).
			WithArgs(amount, uid, model.MinBalance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(model.QueryInsertTransaction)).
			WithArgs(uid, 0, amount, model.TransactionTypeWithdraw).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := walletRepo.Withdraw(ctx, uid, amount)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TestWithdraw_UpdateError", func(t *testing.T) {
		uid := int64(456)
		amount := decimal.NewFromFloat(100.5)
		expectedErr := fmt.Errorf("update failed")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletWithdraw)).
			WithArgs(amount, uid, model.MinBalance).
			WillReturnError(expectedErr)
		mock.ExpectRollback()

		err := walletRepo.Withdraw(ctx, uid, amount)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TestWithdraw_InsertError", func(t *testing.T) {
		uid := int64(789)
		amount := decimal.NewFromFloat(100.5)
		expectedErr := fmt.Errorf("insert failed")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletWithdraw)).
			WithArgs(amount, uid, model.MinBalance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(model.QueryInsertTransaction)).
			WithArgs(uid, 0, amount, model.TransactionTypeWithdraw).
			WillReturnError(expectedErr)
		mock.ExpectRollback()

		err := walletRepo.Withdraw(ctx, uid, amount)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWalletRepo_Transfer(t *testing.T) {
	defer goleak.VerifyNone(t)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	walletRepo := &WalletRepo{
		db:     db,
		logger: zap.NewExample().Sugar(),
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	t.Run("TestTransfer_Success", func(t *testing.T) {
		fromUID := int64(123)
		toUID := int64(456)
		amount := decimal.NewFromFloat(100.5)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletWithdraw)).
			WithArgs(amount, fromUID, model.MinBalance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletTransfer)).
			WithArgs(amount, toUID, model.MaxBalance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(model.QueryInsertTransaction)).
			WithArgs(fromUID, toUID, amount, model.TransactionTypeTransfer).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := walletRepo.Transfer(ctx, fromUID, toUID, amount)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TestTransfer_WithdrawError", func(t *testing.T) {
		fromUID := int64(123)
		toUID := int64(456)
		amount := decimal.NewFromFloat(100.5)
		expectedErr := fmt.Errorf("withdraw failed")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletWithdraw)).
			WithArgs(amount, fromUID, model.MinBalance).
			WillReturnError(expectedErr)
		mock.ExpectRollback()

		err := walletRepo.Transfer(ctx, fromUID, toUID, amount)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TestTransfer_TransferError", func(t *testing.T) {
		fromUID := int64(123)
		toUID := int64(456)
		amount := decimal.NewFromFloat(100.5)
		expectedErr := fmt.Errorf("transfer failed")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletWithdraw)).
			WithArgs(amount, fromUID, model.MinBalance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletTransfer)).
			WithArgs(amount, toUID, model.MaxBalance).
			WillReturnError(expectedErr)
		mock.ExpectRollback()

		err := walletRepo.Transfer(ctx, fromUID, toUID, amount)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TestTransfer_InsertTransactionError", func(t *testing.T) {
		fromUID := int64(123)
		toUID := int64(456)
		amount := decimal.NewFromFloat(100.5)
		expectedErr := fmt.Errorf("insert transaction failed")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletWithdraw)).
			WithArgs(amount, fromUID, model.MinBalance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(model.QueryWalletTransfer)).
			WithArgs(amount, toUID, model.MaxBalance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(model.QueryInsertTransaction)).
			WithArgs(fromUID, toUID, amount, model.TransactionTypeTransfer).
			WillReturnError(expectedErr)
		mock.ExpectRollback()

		err := walletRepo.Transfer(ctx, fromUID, toUID, amount)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWalletRepo_Balance(t *testing.T) {
	defer goleak.VerifyNone(t)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	walletRepo := &WalletRepo{
		db:     db,
		logger: zap.NewExample().Sugar(),
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	t.Run("TestBalance_Normal", func(t *testing.T) {
		uid := int64(123)
		expectedBalance := decimal.NewFromFloat(100.5)

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletBalance)).
			WithArgs(uid).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(expectedBalance))

		balance, errBalance := walletRepo.Balance(ctx, uid)
		require.NoError(t, errBalance)
		assert.Equal(t, expectedBalance, balance)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TestBalance_EmptyResult", func(t *testing.T) {
		uid := int64(456)
		expectedErr := errors.New("sql: no rows in result set")

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletBalance)).
			WithArgs(uid).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}))

		balance, errBalance := walletRepo.Balance(ctx, uid)
		assert.Equal(t, expectedErr, errBalance)
		assert.Equal(t, decimal.Zero, balance)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TestBalance_QueryError", func(t *testing.T) {
		uid := int64(789)
		expectedBalance := decimal.New(0, 1)
		expectedErr := fmt.Errorf("simulated query error")

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryWalletBalance)).
			WithArgs(uid).
			WillReturnError(expectedErr)

		balance, errBalance := walletRepo.Balance(ctx, uid)

		require.Error(t, errBalance)
		assert.Equal(t, expectedBalance, balance)
		assert.Equal(t, expectedErr, errBalance)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
