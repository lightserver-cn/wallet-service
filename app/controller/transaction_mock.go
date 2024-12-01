package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"server/app/request"
)

// MockTransactionInter is a mock implementation of TransactionInter
type MockTransactionInter struct {
	mock.Mock
}

func (m *MockTransactionInter) GetTransactionsByUID(ctx *gin.Context, req *request.ReqTransactions) (*request.ResTransactions, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*request.ResTransactions), args.Error(1)
}
