package service

import (
	"github.com/gin-gonic/gin"

	"server/app/repository"
	"server/app/request"
)

func NewTransaction(repo repository.TransactionInter) TransactionInter {
	return &TransactionServ{
		repo: repo,
	}
}

type TransactionInter interface {
	GetTransactionsByUID(ctx *gin.Context, req *request.ReqTransactions) (*request.ResTransactions, error)
}

type TransactionServ struct {
	repo repository.TransactionInter
}

func (t *TransactionServ) GetTransactionsByUID(ctx *gin.Context,
	req *request.ReqTransactions) (*request.ResTransactions, error) {
	return t.repo.GetTransactionsByUID(ctx, req)
}
