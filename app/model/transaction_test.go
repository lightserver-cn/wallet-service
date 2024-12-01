package model

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"

	"testing"
)

func TestTransactionTypeString(t *testing.T) {
	defer goleak.VerifyNone(t)

	tests := []struct {
		name     string
		tType    TransactionType
		expected string
	}{
		{"Deposit", TransactionTypeDeposit, Deposit},
		{"Withdraw", TransactionTypeWithdraw, Withdraw},
		{"Transfer", TransactionTypeTransfer, Transfer},
		{"Unknown", 4, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTransactionTypeString(tt.tType)
			assert.Equal(t, tt.expected, result)
		})
	}
}
