package model

import (
	"time"
)

// User represents a user in the system.
type User struct {
	ID           int64      `db:"id" json:"id"`
	Username     string     `db:"username" json:"username"`
	Email        string     `db:"email" json:"email"`
	PasswordHash []byte     `db:"password_hash" json:"-"`
	Status       UserStatus `db:"status" json:"status"` // 1-Valid, 2-Invalid, 3-Disabled
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

const TableNameUser = `t_user`

type UserStatus uint8

const (
	_ UserStatus = iota
	UserStatusValid
	UserStatusInvalid
	UserStatusDisabled
)

const FirstColumnUser = `id, username, email, status, created_at, updated_at`

const QueryUserInsert = `INSERT INTO ` + TableNameUser + `(username, email, password_hash) VALUES($1, $2, $3) RETURNING id`
const LogUserInsert = `INSERT INTO ` + TableNameUser + `(username, email, password_hash) VALUES(%s, %s, %s) RETURNING id`

const QueryUserUpdate = `UPDATE ` + TableNameUser + ` SET username=$1, email=$2 WHERE id=$3`
const LogUserUpdate = `UPDATE ` + TableNameUser + ` SET username='%s', email='%s' WHERE id=%d`

const QueryUserBy = `SELECT ` + FirstColumnUser + ` FROM ` + TableNameUser + ` WHERE`
const LogUserByField = QueryUserBy + ` %s = %v`

const QueryUserByID = QueryUserBy + ` id = $1`
const QueryUserByUsername = QueryUserBy + ` username = $1`
const QueryUserByEmail = QueryUserBy + ` email = $1`

var QueryByFieldMap = map[string]string{
	"id":       QueryUserByID,
	"username": QueryUserByUsername,
	"email":    QueryUserByEmail,
}

func GetQueryByField(field string) string {
	str, ok := QueryByFieldMap[field]
	if !ok {
		return ""
	}

	return str
}
