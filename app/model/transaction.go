package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Transaction represents a transaction between wallets.
type Transaction struct {
	ID               int64           `db:"id" json:"id"`
	SenderWalletID   int64           `db:"sender_wallet_id" json:"sender_wallet_id"`     // Foreign key to Wallet.ID
	ReceiverWalletID int64           `db:"receiver_wallet_id" json:"receiver_wallet_id"` // Foreign key to Wallet.ID, can be null
	Amount           decimal.Decimal `db:"amount" json:"amount"`
	TransactionType  TransactionType `db:"transaction_type" json:"transaction_type"` // 1-deposit, 2-withdraw, 3-transfer
	CreatedAt        time.Time       `db:"created_at" json:"created_at"`
}

type TransactionWithUsername struct {
	Transaction
	SenderUsername      string `db:"sender_username" json:"sender_username"`
	ReceiverUsername    string `db:"receiver_username" json:"receiver_username"`
	TransactionTypeName string `json:"transaction_type_name"`
}

const TableNameTransaction = `t_transaction`
const ListColumnTransaction = `t.id, t.sender_wallet_id, COALESCE(s.username, '') AS sender_username, 
		t.receiver_wallet_id, COALESCE(r.username, '') AS receiver_username, amount, t.transaction_type, t.created_at`
const QueryInsertTransaction = `INSERT INTO ` + TableNameTransaction + `
    (sender_wallet_id, receiver_wallet_id, amount, transaction_type, created_at) 
					VALUES ($1, $2, $3, $4, NOW())`
const LogInsertTransaction = `INSERT INTO ` + TableNameTransaction + ` 
    (sender_wallet_id, receiver_wallet_id, amount, transaction_type, created_at) 
					VALUES (%d, %d, %v, %d, NOW())`

const QueryListTransaction = `SELECT ` + ListColumnTransaction + ` FROM ` + TableNameTransaction + ` AS t
		LEFT JOIN t_user AS s ON t.sender_wallet_id = s.id
		LEFT JOIN t_user AS r ON t.receiver_wallet_id = r.id
		WHERE 
 			(t.sender_wallet_id = $1 OR t.receiver_wallet_id = $1) 
			AND (
        		NOT ($2::smallint BETWEEN $3::smallint AND $4::smallint) 
        		OR (t.transaction_type = $2::smallint)
    		)
		ORDER BY t.id DESC
		LIMIT $5 OFFSET $6`
const LogListTransaction = `SELECT ` + ListColumnTransaction + ` FROM ` + TableNameTransaction + ` AS t
		LEFT JOIN t_user AS s ON t.sender_wallet_id = s.id
		LEFT JOIN t_user AS r ON t.receiver_wallet_id = r.id
		WHERE 
			(t.sender_wallet_id = %d OR t.receiver_wallet_id = %d)
			AND (
        		NOT (%d::smallint BETWEEN %d::smallint AND %d::smallint) 
				OR (t.transaction_type = %d::smallint)
    		)
		ORDER BY t.id DESC
		LIMIT %d OFFSET %d`

// TransactionType represents the type of transaction
type TransactionType uint8

const (
	_ TransactionType = iota
	TransactionTypeDeposit
	TransactionTypeWithdraw
	TransactionTypeTransfer
)

const (
	Deposit  = "deposit"
	Withdraw = "withdraw"
	Transfer = "transfer"
)

var transactionTypeMap = map[TransactionType]string{
	TransactionTypeDeposit:  Deposit,
	TransactionTypeWithdraw: Withdraw,
	TransactionTypeTransfer: Transfer,
}

// GetTransactionTypeString returns the string representation of the TransactionType
// If the TransactionType does not exist, it returns an empty string.
func GetTransactionTypeString(tType TransactionType) string {
	str, ok := transactionTypeMap[tType]
	if !ok {
		return ""
	}

	return str
}
