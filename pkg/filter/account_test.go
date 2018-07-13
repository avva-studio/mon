package filter_test

import (
	"testing"
	"time"

	"github.com/glynternet/go-accounting/accountingtest"
	"github.com/glynternet/go-money/currency"
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

func TestCurrency(t *testing.T) {

	for _, test := range []struct {
		name string
		storage.Account
		code  currency.Code
		match bool
	}{
		{
			name:  "zero-values",
			match: true,
		},
		{
			name:    "nil code with valid account",
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "AUP"), time.Time{})},
		},
		{
			name: "nil account with valid code",
			code: accountingtest.NewCurrencyCode(t, "AUP"),
		},
		{
			name:    "valid code and account with non-matching code",
			code:    accountingtest.NewCurrencyCode(t, "BUP"),
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "AUP"), time.Time{})},
		},
		{
			name:    "valid code and account with matching code",
			code:    accountingtest.NewCurrencyCode(t, "BUP"),
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "BUP"), time.Time{})},
			match:   true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := filter.Currency(test.code)
			match := f(test.Account)
			assert.Equal(t, test.match, match)
		})
	}
}

func TestExisted(t *testing.T) {
	for _, test := range []struct {
		name string
		storage.Account
		time.Time
		match bool
	}{
		{
			name:  "zero-values",
			match: true,
		},
		{
			name:    "time before open",
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "BUP"), time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC))},
			Time:    time.Date(1999, 0, 0, 0, 0, 0, 0, time.UTC),
			match:   false,
		},
		{
			name:    "time equal open",
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "BUP"), time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC))},
			Time:    time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
			match:   true,
		},
		{
			name:    "time after open",
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "BUP"), time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC))},
			Time:    time.Date(2001, 0, 0, 0, 0, 0, 0, time.UTC),
			match:   true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := filter.Existed(test.Time)
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
