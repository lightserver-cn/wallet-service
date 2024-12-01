package repository

import (
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"server/app/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"go.uber.org/zap"
)

func TestUserRepo_NewUser(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("TestNewTransaction", func(t *testing.T) {
		db, _, errNew := sqlmock.New()
		require.NoError(t, errNew)
		defer db.Close()

		inter := NewUser(db, zap.NewExample().Sugar())
		assert.NotNil(t, inter)

		repo, ok := inter.(*UserRepo)
		assert.True(t, ok)
		assert.Equal(t, db, repo.db)
	})

	t.Run("TestNewTransaction_NilDB", func(t *testing.T) {
		inter := NewUser(nil, nil)
		expectedInter := &UserRepo{db: nil}
		assert.Equal(t, expectedInter, inter)
	})
}

func TestUserRepo_GetUserByID(t *testing.T) {
	defer goleak.VerifyNone(t)

	db, mock, errNew := sqlmock.New()
	require.NoError(t, errNew)
	defer db.Close()

	userRepo := &UserRepo{
		db:     db,
		logger: zap.NewExample().Sugar(),
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	t.Run("GetUserByID_Normal", func(t *testing.T) {
		id := int64(1)
		expectedUser := &model.User{
			ID:        1,
			Username:  "testuser",
			Email:     "test@example.com",
			Status:    model.UserStatusValid,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryUserByID)).
			WithArgs(id).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "username", "email", "status", "created_at", "updated_at"}).
					AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email, expectedUser.Status,
						expectedUser.CreatedAt, expectedUser.UpdatedAt))

		user, err := userRepo.GetUserByID(ctx, id)

		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetUserByID_NoRows", func(t *testing.T) {
		id := int64(2)
		expectedUser := &model.User{}
		expectedErr := fmt.Errorf("sql: no rows in result set")

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryUserByID)).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "status", "created_at", "updated_at"}))

		user, err := userRepo.GetUserByID(ctx, id)

		assert.Equal(t, expectedUser, user)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetUserByID_PrepareError", func(t *testing.T) {
		id := int64(3)
		expectedUser := &model.User{}

		expectedErr := fmt.Errorf("failed to query model by field: %w", fmt.Errorf("simulated prepare error"))
		mock.ExpectQuery(regexp.QuoteMeta(model.QueryUserByID)).
			WillReturnError(expectedErr)

		user, err := userRepo.GetUserByID(ctx, id)

		assert.Equal(t, expectedUser, user)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetUserByID_QueryError", func(t *testing.T) {
		id := int64(4)
		expectedUser := &model.User{}
		expectedErr := fmt.Errorf("simulated query error")

		mock.ExpectQuery(regexp.QuoteMeta(model.QueryUserByID)).
			WithArgs(id).
			WillReturnError(expectedErr)

		user, err := userRepo.GetUserByID(ctx, id)

		assert.Equal(t, expectedUser, user)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
