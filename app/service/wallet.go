package service

import (
	"fmt"

	"server/app/model"
	"server/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// NewWallet creates a new Wallet service instance.
func NewWallet(repo repository.WalletInter) WalletInter {
	return &WalletServ{
		repo: repo,
	}
}

// WalletInter defines the interface for wallet operations.
type WalletInter interface {
	Deposit(ctx *gin.Context, uid int64, amount decimal.Decimal) error
	Withdraw(ctx *gin.Context, uid int64, amount decimal.Decimal) error
	Transfer(ctx *gin.Context, fromUID, toUID int64, amount decimal.Decimal) error
	Balance(ctx *gin.Context, uid int64) (decimal.Decimal, error)
}

// WalletServ implements the WalletInter interface.
type WalletServ struct {
	repo repository.WalletInter
}

// Deposit adds the specified amount to the user's balance.
func (w *WalletServ) Deposit(ctx *gin.Context, uid int64, amount decimal.Decimal) error {
	// Check if the deposit amount is positive
	if amount.LessThan(decimal.Zero) {
		return fmt.Errorf("deposit amount must be positive")
	}

	// Get the current balance of the user
	balance, err := w.repo.Balance(ctx, uid)
	if err != nil {
		return err
	}

	// Check if the deposit would exceed the maximum allowed balance
	maxBalance := decimal.NewFromInt(model.MaxBalance)
	if balance.Add(amount).GreaterThan(maxBalance) {
		return fmt.Errorf("deposit would exceed the maximum allowed balance of %s", maxBalance.String())
	}

	// Perform the deposit operation
	return w.repo.Deposit(ctx, uid, amount)
}

// Withdraw subtracts the specified amount from the user's balance.
func (w *WalletServ) Withdraw(ctx *gin.Context, uid int64, amount decimal.Decimal) error {
	// Check if the withdraw amount is positive
	if amount.LessThan(decimal.Zero) {
		return fmt.Errorf("withdraw amount must be positive")
	}

	// Get the current balance of the user
	balance, err := w.repo.Balance(ctx, uid)
	if err != nil {
		return err
	}

	// Check if the user has sufficient balance for the withdrawal
	if balance.LessThan(amount) {
		return fmt.Errorf("insufficient balance for withdrawal")
	}

	// Perform the withdrawal operation
	return w.repo.Withdraw(ctx, uid, amount)
}

// Transfer moves the specified amount from the sender's balance to the receiver's balance.
func (w *WalletServ) Transfer(ctx *gin.Context, fromUID, toUID int64, amount decimal.Decimal) error {
	// Check if the transfer amount is positive
	if amount.LessThan(decimal.Zero) {
		return fmt.Errorf("transfer amount must be positive")
	}

	// Get the current balance of the sender
	fromBalance, err := w.repo.Balance(ctx, fromUID)
	if err != nil {
		return err
	}

	// Check if the sender has sufficient balance for the transfer
	if fromBalance.LessThan(amount) {
		return fmt.Errorf("insufficient balance for transfer")
	}

	// Get the current balance of the receiver
	toBalance, err := w.repo.Balance(ctx, toUID)
	if err != nil {
		return err
	}

	// Check if the transfer would exceed the maximum allowed balance for the receiver
	maxBalance := decimal.NewFromInt(model.MaxBalance)
	if toBalance.Add(amount).GreaterThan(maxBalance) {
		return fmt.Errorf("transfer would exceed the maximum allowed balance of %s for the receiver", maxBalance.String())
	}

	// Perform the transfer operation
	return w.repo.Transfer(ctx, fromUID, toUID, amount)
}

// Balance returns the current balance of the user.
func (w *WalletServ) Balance(ctx *gin.Context, uid int64) (decimal.Decimal, error) {
	return w.repo.Balance(ctx, uid)
}
