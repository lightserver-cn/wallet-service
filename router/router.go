package router

import (
	"database/sql"
	"net/http"

	"go.uber.org/zap"

	"server/app/controller"
	"server/app/repository"
	"server/app/request"
	"server/app/service"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine, db *sql.DB, logger *zap.SugaredLogger) {
	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, request.ResponseEntity{
			ErrCode: 0,
			ErrMsg:  "success",
			Data:    struct{}{},
		})
	})

	userRepo := repository.NewUser(db, logger)
	walletRepo := repository.NewWallet(db, logger)
	transactionRepo := repository.NewTransaction(db, logger)

	userServ := service.NewUser(userRepo, walletRepo)
	userCtrl := controller.NewUser(userServ)

	userRout := router.Group("/api/users")
	userRout.POST("", userCtrl.RegisterUser)
	userRout.GET("/:uid", userCtrl.GetUserByUID)

	transactionServ := service.NewTransaction(transactionRepo)
	walletServ := service.NewWallet(walletRepo)
	walletCtrl := controller.NewWallet(walletServ, transactionServ)

	walletRout := router.Group("/api/wallets")
	walletRout.POST("/:uid/deposit", walletCtrl.Deposit)
	walletRout.POST("/:uid/withdraw", walletCtrl.Withdraw)
	walletRout.POST("/:uid/transfer", walletCtrl.Transfer)
	walletRout.GET("/:uid/balance", walletCtrl.Balance)
	walletRout.GET("/:uid/transactions", walletCtrl.Transactions)
}
