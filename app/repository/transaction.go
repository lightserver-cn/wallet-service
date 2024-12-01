package repository

import (
	"database/sql"
	"fmt"

	"server/app/model"
	"server/app/request"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewTransaction(db *sql.DB, logger *zap.SugaredLogger) TransactionInter {
	return &TransactionRepo{
		db:     db,
		logger: logger,
	}
}

type TransactionInter interface {
	GetTransactionsByUID(ctx *gin.Context, req *request.ReqTransactions) (*request.ResTransactions, error)
}

type TransactionRepo struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

// GetTransactionsByUID retrieves a list of transactions related to a user ID with pagination.
func (t *TransactionRepo) GetTransactionsByUID(ctx *gin.Context,
	req *request.ReqTransactions) (*request.ResTransactions, error) {
	res := &request.ResTransactions{}

	req.ValidatePageSize()
	// Calculate offset based on page number and page size
	offset := (req.Page - 1) * req.PageSize

	t.logger.Infof(model.LogListTransaction, req.UID, req.UID, req.Type,
		model.TransactionTypeDeposit, model.TransactionTypeTransfer, req.Type, req.PageSize+1, offset)

	rows, err := t.db.QueryContext(ctx, model.QueryListTransaction, req.UID, req.Type,
		model.TransactionTypeDeposit, model.TransactionTypeTransfer, req.PageSize+1, offset)
	if err != nil {
		t.logger.Errorf("query transactions error: %s", err)
		return res, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var transactions []*model.TransactionWithUsername
	for rows.Next() {
		mod := &model.TransactionWithUsername{}

		err = rows.Scan(&mod.ID, &mod.SenderWalletID, &mod.SenderUsername, &mod.ReceiverWalletID, &mod.ReceiverUsername,
			&mod.Amount, &mod.TransactionType, &mod.CreatedAt)
		if err != nil {
			t.logger.Errorf("GetTransactionsByUID failed to scan rows: %v", err)
			return res, fmt.Errorf("failed to scan row: %w", err)
		}

		mod.TransactionTypeName = model.GetTransactionTypeString(mod.TransactionType)

		transactions = append(transactions, mod)
	}

	if rows.Err() != nil {
		t.logger.Errorf("GetTransactionsByUID failed to scan rows: %v", rows.Err())
		return res, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	hasMore := len(transactions) == req.PageSize+1
	if hasMore {
		transactions = transactions[:req.PageSize]
	}

	res.List = transactions
	res.HasMore = hasMore

	return res, nil
}
