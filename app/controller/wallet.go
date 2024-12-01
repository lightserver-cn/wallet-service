package controller

import (
	"database/sql"
	"errors"
	"net/http"

	"server/app/model"
	"server/app/request"
	"server/app/service"
	"server/pkg/consts"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func NewWallet(serv service.WalletInter, servTransaction service.TransactionInter) WalletInter {
	return &WalletCtrl{
		serv:            serv,
		servTransaction: servTransaction,
	}
}

type WalletInter interface {
	Deposit(ctx *gin.Context)
	Withdraw(ctx *gin.Context)
	Transfer(ctx *gin.Context)
	Balance(ctx *gin.Context)
	Transactions(ctx *gin.Context)
}

type WalletCtrl struct {
	serv            service.WalletInter
	servTransaction service.TransactionInter
}

// handleWalletOperation is a generic handler function used to process deposit and withdrawal operations.
func handleWalletOperation(ctx *gin.Context, operation func(ctx *gin.Context, uid int64, amount decimal.Decimal) error) {
	idReq := new(request.ReqUID)
	if err := ctx.ShouldBindUri(idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrValidationFailed, "details": err.Error()})
		return
	}

	if idReq.UID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrInvalidUID})
		return
	}

	amountReq := new(request.ReqAmount)
	if err := ctx.ShouldBindJSON(amountReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrValidationFailed, "details": err.Error()})
		return
	}

	if amountReq.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrInvalidAmount})
		return
	}

	err := operation(ctx, idReq.UID, amountReq.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": consts.ErrInternalServer, "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": consts.MsgSuccess})
}

func (w *WalletCtrl) Deposit(ctx *gin.Context) {
	handleWalletOperation(ctx, w.serv.Deposit)
}

func (w *WalletCtrl) Withdraw(ctx *gin.Context) {
	handleWalletOperation(ctx, w.serv.Withdraw)
}

func (w *WalletCtrl) Transfer(ctx *gin.Context) {
	idReq := new(request.ReqUID)
	if err := ctx.ShouldBindUri(idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrValidationFailed, "details": err.Error()})
		return
	}

	if idReq.UID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrInvalidUID})
		return
	}

	transferReq := new(request.ReqTransfer)
	if err := ctx.ShouldBindJSON(transferReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrValidationFailed, "details": err.Error()})
		return
	}

	if transferReq.ToUID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrInvalidUID})
		return
	}

	if transferReq.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrInvalidAmount})
		return
	}

	err := w.serv.Transfer(ctx, idReq.UID, transferReq.ToUID, transferReq.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": consts.ErrTransferFailed, "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": consts.MsgSuccess})
}

func (w *WalletCtrl) Balance(ctx *gin.Context) {
	var idReq request.ReqUID
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   consts.ErrValidationFailed,
			"details": err.Error(),
		})
		return
	}

	if idReq.UID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrInvalidUID})
		return
	}

	balance, err := w.serv.Balance(ctx, idReq.UID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": consts.ErrUserNotFound})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": consts.ErrInternalServer, "details": err.Error()})
		}
		return
	}

	res := &request.ResBalance{
		Balance: balance,
	}

	ctx.JSON(http.StatusOK, res)
}

func (w *WalletCtrl) Transactions(ctx *gin.Context) {
	idReq := new(request.ReqUID)
	if err := ctx.ShouldBindUri(idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrValidationFailed, "details": err.Error()})
		return
	}

	if idReq.UID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrInvalidUID})
		return
	}

	req := new(request.ReqTransactions)
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrValidationFailed, "details": err.Error()})
		return
	}
	req.ValidatePageSize()

	if req.Type > model.TransactionTypeTransfer {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrInvalidTransactionType})
		return
	}

	req.UID = idReq.UID

	res, err := w.servTransaction.GetTransactionsByUID(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": consts.ErrInternalServer, "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
