package consts

import "time"

const TimestampFormat1 = "2006-01-02 15:04:05"

const (
	Day   = 24 * time.Hour
	Day30 = 30 * Day
)

const Thousand = 1000

const (
	MsgSuccess                = "Successful"
	ErrValidationFailed       = "Validation failed"
	ErrUsernameRequired       = "username is required"
	ErrEmailRequired          = "email is required"
	ErrPasswordRequired       = "password is required"
	ErrUsernameAlreadyExists  = "The username has been repeated, please change to a new username."
	ErrEmailAlreadyExists     = "The email has been repeated, please change to a new email."
	ErrInternalServer         = "internal server error"
	ErrInvalidUID             = "Invalid UID"
	ErrUserNotFound           = "user not found"
	ErrInvalidAmount          = "Invalid Amount"
	ErrTransferFailed         = "Transfer failed"
	ErrInvalidTransactionType = "Invalid transaction type"
)
