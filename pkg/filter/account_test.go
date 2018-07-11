package filter

import (
	"testing"

	"github.com/glynternet/mon/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func stubFilter(result bool) AccountFilter {
	return func(_ storage.Account) bool {
		return result
	}
}

func TestAccountFilters_Or(t *testing.T) {
	for _, test := range []struct {
		name string
		AccountFilters
		storage.Account
		expected bool
	}{
		{
			name: "zero-values",
		},
		{
			name:     "single filter passing",
			expected: true,
			AccountFilters: AccountFilters{
				stubFilter(true),
			},
		},
		{
			name: "single filter failing",
			AccountFilters: AccountFilters{
				stubFilter(false),
			},
		},
		{
			name:     "multiple filters passing",
			expected: true,
			AccountFilters: AccountFilters{
				stubFilter(true),
				stubFilter(true),
			},
		},
		{
			name: "multiple filters failing",
			AccountFilters: AccountFilters{
				stubFilter(false),
				stubFilter(false),
			},
		},
		{
			name:     "multiple filters mixed",
			expected: true,
			AccountFilters: AccountFilters{
				stubFilter(false),
				stubFilter(false),
				stubFilter(true),
				stubFilter(false),
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			actual := test.AccountFilters.Or(test.Account)
			assert.Equal(t, test.expected, actual)
		})
	}
}
