package repository

import (
	"errors"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"server/app/model"
	"server/app/request"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"go.uber.org/zap"
)

func TestTransactionRepo_NewTransaction(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("TestNewTransaction", func(t *testing.T) {
		db, _, errNew := sqlmock.New()
		require.NoError(t, errNew)
		defer db.Close()

		inter := NewTransaction(db, zap.NewExample().Sugar())
		assert.NotNil(t, inter)

		repo, ok := inter.(*TransactionRepo)
		assert.True(t, ok)
		assert.Equal(t, db, repo.db)
	})

	t.Run("TestNewTransaction_NilDB", func(t *testing.T) {
		inter := NewTransaction(nil, nil)
		expectedInter := &TransactionRepo{db: nil}
		assert.Equal(t, expectedInter, inter)
	})
}

func TestGetTransactionsByUID(t *testing.T) {
	defer goleak.VerifyNone(t)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &TransactionRepo{
		db:     db,
		logger: zap.NewExample().Sugar(),
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	columns := []string{
		"id", "sender_wallet_id", "sender_username", "receiver_wallet_id", "receiver_username",
		"amount", "transaction_type", "created_at",
	}

	req := &request.ReqTransactions{
		UID: 1,
		ReqPage: request.ReqPage{
			Page:     1,
			PageSize: 10,
		},
	}

	t.Run("Test with invalid page number", func(t *testing.T) {
		req.Page = -1 // Invalid page number

		rows := sqlmock.NewRows([]string{})
		mock.ExpectQuery(regexp.QuoteMeta(model.QueryListTransaction)).
			WithArgs(req.UID, req.Type, model.TransactionTypeDeposit, model.TransactionTypeTransfer, req.PageSize+1, 0).
			WillReturnRows(rows)

		res, err := repo.GetTransactionsByUID(ctx, req)
		require.NoError(t, err)
		assert.Empty(t, res.List)
		assert.False(t, res.HasMore)
	})

	t.Run("Test with invalid page size", func(t *testing.T) {
		req.PageSize = -1 // Invalid page size

		rows := sqlmock.NewRows([]string{})

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryListTransaction)).
			WithArgs(req.UID, req.Type, model.TransactionTypeDeposit, model.TransactionTypeTransfer, 1+1, 0).
			WillReturnRows(rows)

		res, err := repo.GetTransactionsByUID(ctx, req)
		require.NoError(t, err)
		assert.Empty(t, res.List)
		assert.False(t, res.HasMore)
	})

	t.Run("Test with large page size", func(t *testing.T) {
		req.PageSize = 10000 // Large page size

		rows := sqlmock.NewRows([]string{})

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryListTransaction)).
			WithArgs(req.UID, req.Type, model.TransactionTypeDeposit, model.TransactionTypeTransfer, 100+1, 0).
			WillReturnRows(rows)

		res, err := repo.GetTransactionsByUID(ctx, req)
		require.NoError(t, err)
		assert.Empty(t, res.List)
		assert.False(t, res.HasMore)
	})

	t.Run("Test with empty result set", func(t *testing.T) {
		rows := sqlmock.NewRows(columns)

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryListTransaction)).
			WithArgs(req.UID, req.Type, model.TransactionTypeDeposit, model.TransactionTypeTransfer, req.PageSize+1, 0).
			WillReturnRows(rows)

		res, err := repo.GetTransactionsByUID(ctx, req)
		require.NoError(t, err)
		assert.Empty(t, res.List)
		assert.False(t, res.HasMore)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test with valid input", func(t *testing.T) {
		rows := sqlmock.NewRows(columns).
			AddRow(1, 101, "sender1", 102, "receiver1", 100.0, model.TransactionTypeDeposit, time.Now()).
			AddRow(2, 103, "sender2", 104, "receiver2", 200.0, model.TransactionTypeWithdraw, time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryListTransaction)).
			WithArgs(req.UID, req.Type, model.TransactionTypeDeposit, model.TransactionTypeTransfer, req.PageSize+1, 0).
			WillReturnRows(rows)

		res, err := repo.GetTransactionsByUID(ctx, req)
		require.NoError(t, err)
		assert.Len(t, res.List, 2)
		assert.False(t, res.HasMore)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test with no transactions", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{})

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryListTransaction)).
			WithArgs(req.UID, req.Type, model.TransactionTypeDeposit, model.TransactionTypeTransfer, req.PageSize+1, 0).
			WillReturnRows(rows)

		res, err := repo.GetTransactionsByUID(ctx, req)
		require.NoError(t, err)
		assert.Empty(t, res.List)
		assert.False(t, res.HasMore)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test with error preparing statement", func(t *testing.T) {
		expectedRes := &request.ResTransactions{
			List:    []*model.TransactionWithUsername(nil),
			HasMore: false,
		}

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryListTransaction)).
			WithArgs(req.UID, req.Type, model.TransactionTypeDeposit, model.TransactionTypeTransfer, req.PageSize+1, 0).
			WillReturnError(errors.New("statement preparation error"))

		res, err := repo.GetTransactionsByUID(ctx, req)
		require.Error(t, err)
		assert.Equal(t, expectedRes, res)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test with error executing query", func(t *testing.T) {
		expectedRes := &request.ResTransactions{
			List:    []*model.TransactionWithUsername(nil),
			HasMore: false,
		}

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryListTransaction)).
			WithArgs(req.UID, req.Type, model.TransactionTypeDeposit, model.TransactionTypeTransfer, req.PageSize+1, 0).
			WillReturnError(errors.New("query execution error"))

		res, err := repo.GetTransactionsByUID(ctx, req)
		require.Error(t, err)
		assert.Equal(t, expectedRes, res)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test with error scanning row", func(t *testing.T) {
		rows := sqlmock.NewRows(columns).
			AddRow(1, 101, "sender1", 102, "receiver1", 100.0, model.TransactionTypeWithdraw, "2023-04-01").
			AddRow(2, 103, "sender2", 104, "receiver2", 200.0, model.TransactionTypeDeposit, "invalid date")

		expectedRes := &request.ResTransactions{
			List:    []*model.TransactionWithUsername(nil),
			HasMore: false,
		}

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryListTransaction)).
			WithArgs(req.UID, req.Type, model.TransactionTypeDeposit, model.TransactionTypeTransfer, req.PageSize+1, 0).
			WillReturnRows(rows)

		res, err := repo.GetTransactionsByUID(ctx, req)
		require.Error(t, err)
		assert.Equal(t, expectedRes, res)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No more data", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{})

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryListTransaction)).
			WithArgs(req.UID, req.Type, model.TransactionTypeDeposit, model.TransactionTypeTransfer, req.PageSize+1, 0).
			WillReturnRows(rows)

		res, err := repo.GetTransactionsByUID(ctx, req)
		require.NoError(t, err)
		assert.Empty(t, res.List)
		assert.False(t, res.HasMore)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
