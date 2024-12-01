package service

import (
	"net/http/httptest"
	"testing"

	"server/app/model"
	"server/app/request"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestTransactionServ_NewTransaction(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("TestNewTransaction", func(t *testing.T) {
		// Create a mock instance
		repo := new(MockTransactionInter)

		inter := NewTransaction(repo)
		assert.NotNil(t, inter)

		serv, ok := inter.(*TransactionServ)
		assert.True(t, ok)
		assert.Equal(t, repo, serv.repo)
	})

	t.Run("TestNewTransaction_NilRepo", func(t *testing.T) {
		inter := NewTransaction(nil)
		expectedInter := &TransactionServ{repo: nil}
		assert.Equal(t, expectedInter, inter)
	})
}

func TestTransactionServ_GetTransactionsByUID(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Create a mock instance
	mockRepo := new(MockTransactionInter)

	// Setup the expected call and return value
	expectedReq := &request.ReqTransactions{
		UID: 1,
		ReqPage: request.ReqPage{
			Page:     1,
			PageSize: 10,
		},
	}

	expectedRes := &request.ResTransactions{
		List:    []*model.TransactionWithUsername{},
		HasMore: false,
	}
	mockRepo.On("GetTransactionsByUID", mock.Anything, expectedReq).Return(expectedRes, nil)

	// Create a TransactionServ instance with the mock repo
	serv := NewTransaction(mockRepo)

	// Create a Gin context
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Call the method under test
	res, err := serv.GetTransactionsByUID(ctx, expectedReq)

	// Assert the results
	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)

	// Assert that the mock was called as expected
	mockRepo.AssertExpectations(t)
}
