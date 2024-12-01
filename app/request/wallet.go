package request

import (
	"github.com/shopspring/decimal"

	"server/app/model"
)

type ReqAmount struct {
	Amount decimal.Decimal `json:"amount"`
}

type ReqDeposit struct {
	Amount decimal.Decimal `json:"amount"`
}

type ReqWithdraw struct {
	Amount decimal.Decimal `json:"amount"`
}

type ReqTransfer struct {
	ToUID  int64           `json:"to_uid"`
	Amount decimal.Decimal `json:"amount"`
}

type ResBalance struct {
	Balance decimal.Decimal `json:"balance"`
}

type ReqTransactions struct {
	UID  int64                 `json:"-"`
	Type model.TransactionType `form:"type" `
	ReqPage
}

type ResTransactions struct {
	List    []*model.TransactionWithUsername `json:"list"`
	HasMore bool                             `json:"has_more"`
}
