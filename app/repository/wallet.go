package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"server/app/model"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func NewWallet(db *sql.DB, logger *zap.SugaredLogger) WalletInter {
	return &WalletRepo{
		db:     db,
		logger: logger,
	}
}

type WalletInter interface {
	CreateWallet(ctx *gin.Context, mod *model.Wallet) (*model.Wallet, error)
	GetWalletByUID(ctx *gin.Context, uid int64) (*model.Wallet, error)
	Deposit(ctx *gin.Context, uid int64, amount decimal.Decimal) error
	Withdraw(ctx *gin.Context, uid int64, amount decimal.Decimal) error
	Transfer(ctx *gin.Context, fromUID, toUID int64, amount decimal.Decimal) error
	Balance(ctx *gin.Context, uid int64) (decimal.Decimal, error)
}

type WalletRepo struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

// CreateWallet create wallet
func (w *WalletRepo) CreateWallet(ctx *gin.Context, mod *model.Wallet) (*model.Wallet, error) {
	var id int64

	w.logger.Infof(model.LogWalletInert, mod.UID, mod.Balance)
	err := w.db.QueryRowContext(ctx, model.QueryWalletInsert, mod.UID, mod.Balance).Scan(&id)
	if err != nil {
		return mod, fmt.Errorf("failed to insert wallet: %w", err)
	}

	mod.ID = id

	w.logger.Infof("Wallet created with ID: %d", id)

	return mod, err
}

// Deposit adds money to a user's wallet and records the transaction
func (w *WalletRepo) Deposit(ctx *gin.Context, uid int64, amount decimal.Decimal) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Errorf("Deposit failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // if Commit returns error update err with commit err
		}
	}()

	w.logger.Infof(model.LogWalletDeposit, amount, uid, amount, model.MaxBalance)
	w.logger.Infof(model.LogInsertTransaction, 0, uid, amount, model.TransactionTypeDeposit)

	_, err = tx.ExecContext(ctx, model.QueryWalletDeposit, amount, uid, model.MaxBalance)
	if err != nil {
		w.logger.Errorf("Deposit failed to query wallet deposit: %v", err)
		return err
	}

	_, err = tx.ExecContext(ctx, model.QueryInsertTransaction, 0, uid, amount, model.TransactionTypeDeposit)
	if err != nil {
		w.logger.Errorf("Deposit failed to query insert transaction: %v", err)
		return err
	}

	return nil
}

// Withdraw removes money from a user's wallet and records the transaction
func (w *WalletRepo) Withdraw(ctx *gin.Context, uid int64, amount decimal.Decimal) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Errorf("Withdraw failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // if Commit returns error update err with commit err
		}
	}()

	w.logger.Infof(model.LogWalletWithdraw, amount, uid, amount, model.MinBalance)
	w.logger.Infof(model.LogInsertTransaction, uid, 0, amount, model.TransactionTypeWithdraw)

	_, err = tx.ExecContext(ctx, model.QueryWalletWithdraw, amount, uid, model.MinBalance)
	if err != nil {
		w.logger.Errorf("Withdraw failed to query wallet withdraw: %v", err)
		return err
	}

	_, err = tx.ExecContext(ctx, model.QueryInsertTransaction, uid, 0, amount, model.TransactionTypeWithdraw)
	if err != nil {
		w.logger.Errorf("Withdraw failed to query insert transaction: %v", err)
		return err
	}

	return nil
}

func (w *WalletRepo) Transfer(ctx *gin.Context, fromUID, toUID int64, amount decimal.Decimal) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Errorf("Transfer failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // if Commit returns error update err with commit err
		}
	}()

	w.logger.Infof(model.LogWalletWithdraw, amount, fromUID, amount, model.MinBalance)
	w.logger.Infof(model.LogWalletTransfer, amount, toUID, amount, model.MaxBalance)
	w.logger.Infof(model.LogInsertTransaction, fromUID, toUID, amount, model.TransactionTypeTransfer)

	_, err = tx.ExecContext(ctx, model.QueryWalletWithdraw, amount, fromUID, model.MinBalance)
	if err != nil {
		_ = tx.Rollback()
		w.logger.Errorf("Transfer failed to query wallet withdraw: %v", err)
		return err
	}

	_, err = tx.ExecContext(ctx, model.QueryWalletTransfer, amount, toUID, model.MaxBalance)
	if err != nil {
		_ = tx.Rollback()
		w.logger.Errorf("Transfer failed to query wallet transfer: %v", err)
		return err
	}

	_, err = tx.ExecContext(ctx, model.QueryInsertTransaction, fromUID, toUID, amount, model.TransactionTypeTransfer)
	if err != nil {
		_ = tx.Rollback()
		w.logger.Errorf("Transfer failed to query inert transaction: %v", err)
		return err
	}

	return tx.Commit()
}

func (w *WalletRepo) Balance(ctx *gin.Context, uid int64) (decimal.Decimal, error) {
	w.logger.Infof(model.LogWalletBalance, uid)

	var balance decimal.Decimal
	err := w.db.QueryRowContext(ctx, model.QueryWalletBalance, uid).Scan(&balance)
	if err != nil {
		w.logger.Errorf("Balance failed to get query wallet balance: %v", err)
		return decimal.Zero, err
	}

	return balance, nil
}

func (w *WalletRepo) GetWalletByUID(ctx *gin.Context, uid int64) (*model.Wallet, error) {
	return w.queryModelByField(ctx, "uid", uid)
}

func (w *WalletRepo) queryModelByField(ctx *gin.Context, field string, value any) (*model.Wallet, error) {
	mod := &model.Wallet{}

	w.logger.Infof(model.LogWalletByField, field, value)

	err := w.db.QueryRowContext(ctx, model.QueryWalletByField, field, value).
		Scan(&mod.ID, &mod.UID, &mod.Balance, &mod.CreatedAt, &mod.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mod, err
		}

		w.logger.Errorf("queryModelByField failed to query wallet by field: %v", err)
		return mod, err
	}

	return mod, nil
}
