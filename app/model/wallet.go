package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Wallet represents a wallet belonging to a user.
type Wallet struct {
	ID        int64           `db:"id" json:"id"`
	UID       int64           `db:"uid" json:"uid"` // Foreign key to User.ID
	Balance   decimal.Decimal `db:"balance" json:"balance"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

const TableNameWallet = `t_wallet`

const (
	MinBalance = 0
	MaxBalance = 1000000
)

const FirstColumnWallet = `id, uid, balance, created_at, updated_at`

const QueryWalletByField = `SELECT ` + FirstColumnWallet + ` FROM ` + TableNameWallet + ` WHERE $2 = $3`
const LogWalletByField = `SELECT ` + FirstColumnWallet + ` FROM ` + TableNameWallet + ` WHERE %s = %v`

const QueryWalletBalance = `SELECT balance FROM ` + TableNameWallet + ` WHERE uid = $1`
const LogWalletBalance = `SELECT balance FROM ` + TableNameWallet + ` WHERE uid = %d`

const QueryWalletInsert = `INSERT INTO ` + TableNameWallet + ` (uid, balance) VALUES($1, $2) RETURNING id`
const LogWalletInert = `INSERT INTO ` + TableNameWallet + ` (uid, balance) VALUES(%d, %v) RETURNING id`

const QueryWalletDeposit = `UPDATE ` + TableNameWallet +
	` SET balance = balance + $1, updated_at = NOW() WHERE uid = $2 AND balance + $1 <= $3`
const LogWalletDeposit = `UPDATE ` + TableNameWallet +
	` SET balance = balance + %v, updated_at = NOW() WHERE uid = %d AND balance + %v <= %d`

const QueryWalletWithdraw = `UPDATE ` + TableNameWallet + ` SET balance = balance - $1, updated_at = NOW() 
		WHERE uid = $2 AND balance - $1 >= $3`
const LogWalletWithdraw = `UPDATE ` + TableNameWallet + ` SET balance = balance - %v, updated_at = NOW() 
		WHERE uid = %d AND balance - %v >= %d`

const QueryWalletTransfer = `UPDATE ` + TableNameWallet + ` SET balance = balance + $1, updated_at = NOW() 
			WHERE uid = $2 AND balance + $1 < $3`
const LogWalletTransfer = `UPDATE ` + TableNameWallet + ` SET balance = balance + %v, updated_at = NOW() 
			WHERE uid = %d AND balance + %v < %d`
