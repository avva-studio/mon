package filter_test

import (
	"testing"

	"github.com/glynternet/mon/pkg/filter"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func stubFilter(result bool) filter.AccountFilter {
	return func(_ storage.Account) bool {
		return result
	}
}

func TestID(t *testing.T) {
	for _, test := range []struct {
		name string
		storage.Account
		id    uint
		match bool
	}{
		{
			name:  "zero-values",
			match: true,
		},
		{
			name:    "matching",
			Account: storage.Account{ID: 111},
			id:      111,
			match:   true,
		},
		{
			name:    "non-matching",
			Account: storage.Account{ID: 222},
			id:      123,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := filter.ID(test.id)
			match := f(test.Account)
			assert.Equal(t, test.match, match)
		})
	}
}

func TestAccountFilters_Or(t *testing.T) {
	for _, test := range []struct {
		name string
		filter.AccountFilters
		storage.Account
		expected bool
	}{
		{
			name: "zero-values",
		},
		{
			name:     "single filter passing",
			expected: true,
			AccountFilters: filter.AccountFilters{
				stubFilter(true),
			},
		},
		{
			name: "single filter failing",
			AccountFilters: filter.AccountFilters{
				stubFilter(false),
			},
		},
		{
			name:     "multiple filters passing",
			expected: true,
			AccountFilters: filter.AccountFilters{
				stubFilter(true),
				stubFilter(true),
			},
		},
		{
			name: "multiple filters failing",
			AccountFilters: filter.AccountFilters{
				stubFilter(false),
				stubFilter(false),
			},
		},
		{
			name:     "multiple filters mixed",
			expected: true,
			AccountFilters: filter.AccountFilters{
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
