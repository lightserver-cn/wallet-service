package request

import (
	"testing"

	"go.uber.org/goleak"
)

func TestValidatePageSize(t *testing.T) {
	defer goleak.VerifyNone(t)

	tests := []struct {
		name     string
		reqPage  ReqPage
		expected ReqPage
	}{
		{
			name: "Valid Page and PageSize",
			reqPage: ReqPage{
				Page:     2,
				PageSize: 10,
			},
			expected: ReqPage{
				Page:     2,
				PageSize: 10,
			},
		},
		{
			name: "Negative Page",
			reqPage: ReqPage{
				Page:     -1,
				PageSize: 10,
			},
			expected: ReqPage{
				Page:     1,
				PageSize: 10,
			},
		},
		{
			name: "Zero Page",
			reqPage: ReqPage{
				Page:     0,
				PageSize: 10,
			},
			expected: ReqPage{
				Page:     1,
				PageSize: 10,
			},
		},
		{
			name: "PageSize less than min",
			reqPage: ReqPage{
				Page:     2,
				PageSize: 0,
			},
			expected: ReqPage{
				Page:     2,
				PageSize: pageSizeMin,
			},
		},
		{
			name: "PageSize greater than max",
			reqPage: ReqPage{
				Page:     2,
				PageSize: 150,
			},
			expected: ReqPage{
				Page:     2,
				PageSize: pageSizeMax,
			},
		},
		{
			name: "Negative PageSize",
			reqPage: ReqPage{
				Page:     2,
				PageSize: -5,
			},
			expected: ReqPage{
				Page:     2,
				PageSize: pageSizeMin,
			},
		},
		{
			name: "All invalid values",
			reqPage: ReqPage{
				Page:     -1,
				PageSize: -5,
			},
			expected: ReqPage{
				Page:     1,
				PageSize: pageSizeMin,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.reqPage.ValidatePageSize()
			if tt.reqPage.Page != tt.expected.Page || tt.reqPage.PageSize != tt.expected.PageSize {
				t.Errorf("ValidatePageSize() = %v, want %v", tt.reqPage, tt.expected)
			}
		})
	}
}
