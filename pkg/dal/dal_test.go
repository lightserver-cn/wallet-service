package dal

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// Mock database connection
	mockDB, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer mockDB.Close()

	// Mock Redis client
	mockRDB := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Mock logger using log.New
	mockLogger := log.New(log.Writer(), "[TEST] ", log.LstdFlags|log.Lshortfile)

	// Test case: Valid parameters
	t.Run("Valid Parameters", func(t *testing.T) {
		dal, err := New(mockDB, mockRDB, mockLogger)
		require.NoError(t, err)
		assert.NotNil(t, dal)
		assert.Equal(t, mockDB, dal.DB)
		assert.Equal(t, mockRDB, dal.RDB)
		assert.Equal(t, mockLogger, dal.logger)
	})

	// Test case: Nil database connection
	t.Run("Nil Database Connection", func(t *testing.T) {
		dal, err := New(nil, mockRDB, mockLogger)
		require.NoError(t, err)
		assert.NotNil(t, dal)
		assert.Nil(t, dal.DB)
		assert.Equal(t, mockRDB, dal.RDB)
		assert.Equal(t, mockLogger, dal.logger)
	})

	// Test case: Nil Redis client
	t.Run("Nil Redis Client", func(t *testing.T) {
		dal, err := New(mockDB, nil, mockLogger)
		require.NoError(t, err)
		assert.NotNil(t, dal)
		assert.Equal(t, mockDB, dal.DB)
		assert.Nil(t, dal.RDB)
		assert.Equal(t, mockLogger, dal.logger)
	})

	// Test case: Nil logger
	t.Run("Nil Logger", func(t *testing.T) {
		dal, err := New(mockDB, mockRDB, nil)
		require.NoError(t, err)
		assert.NotNil(t, dal)
		assert.Equal(t, mockDB, dal.DB)
		assert.Equal(t, mockRDB, dal.RDB)
		assert.Nil(t, dal.logger)
	})
}
