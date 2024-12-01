package model

import (
	"testing"

	"go.uber.org/goleak"
)

func TestGetQueryByField(t *testing.T) {
	defer goleak.VerifyNone(t)

	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{"ByID", "id", QueryUserByID},
		{"ByUsername", "username", QueryUserByUsername},
		{"ByEmail", "email", QueryUserByEmail},
		{"UnknownField", "unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetQueryByField(tt.field)
			if result != tt.expected {
				t.Errorf("GetQueryByField(%q) = %v; want %v", tt.field, result, tt.expected)
			}
		})
	}
}
