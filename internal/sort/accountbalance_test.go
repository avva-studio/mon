package sort_test

import (
	"testing"
	"time"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/common"
	"github.com/glynternet/mon/internal/accountbalance"
	"github.com/glynternet/mon/internal/sort"
	"github.com/stretchr/testify/assert"
)

func TestBalanceAmount(t *testing.T) {
	t.Run("nil input remains nil", func(t *testing.T) {
		var abs []accountbalance.AccountBalance
		sort.BalanceAmount(abs)
		assert.Nil(t, abs)
	})

	accountFn := func(amount int) accountbalance.AccountBalance {
		b, err := balance.New(time.Time{}, balance.Amount(amount))
		common.FatalIfError(t, err, "creating new balance")
		return accountbalance.AccountBalance{Balance: *b}
	}

	for _, test := range []struct {
		name    string
		in, out []int
	}{
		{
			name: "empty accountbalance.AccountBalance",
			in:   []int{},
		},
		{
			name: "single AccountBalance yields same AccountBalance",
			in:   []int{1},
			out:  []int{1},
		},
		{
			name: "in-order input yields same order output",
			in:   []int{-1, 1, 2, 3},
			out:  []int{-1, 1, 2, 3},
		},
		{
			name: "in reverse order",
			in:   []int{3, 2, 1, 0},
			out:  []int{0, 1, 2, 3},
		},
		{
			name: "mixed order",
			in:   []int{3, 7, 99, 2, 5, 8},
			out:  []int{2, 3, 5, 7, 8, 99},
		},
		{
			name: "mixed order with repeating",
			in:   []int{9, 5, 2, 66, 4, 5},
			out:  []int{2, 4, 5, 5, 9, 66},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var abs []accountbalance.AccountBalance
			for _, amount := range test.in {
				abs = append(abs, accountFn(amount))
			}

			sort.BalanceAmount(abs)

			var out []int
			for _, ab := range abs {
				out = append(out, ab.Balance.Amount)
			}

			assert.Equal(t, test.out, out)
		})
	}
}

func TestBalanceMagnitude(t *testing.T) {
	t.Run("nil input remains nil", func(t *testing.T) {
		var abs []accountbalance.AccountBalance
		sort.BalanceAmountMagnitude(abs)
		assert.Nil(t, abs)
	})

	accountFn := func(amount int) accountbalance.AccountBalance {
		b, err := balance.New(time.Time{}, balance.Amount(amount))
		common.FatalIfError(t, err, "creating new balance")
		return accountbalance.AccountBalance{Balance: *b}
	}

	for _, test := range []struct {
		name    string
		in, out []int
	}{
		{
			name: "empty accountbalance.AccountBalance",
			in:   []int{},
		},
		{
			name: "single AccountBalance yields same AccountBalance",
			in:   []int{1},
			out:  []int{1},
		},
		{
			name: "in-order input yields same order output",
			in:   []int{-1, 2, 3},
			out:  []int{-1, 2, 3},
		},
		{
			name: "in reverse order",
			in:   []int{-3, -2, -1, 0},
			out:  []int{0, -1, -2, -3},
		},
		{
			name: "mixed order",
			in:   []int{-3, 7, 99, 2, 5, 8},
			out:  []int{2, -3, 5, 7, 8, 99},
		},
		{
			name: "mixed order with repeating",
			in:   []int{-9, -5, 2, -66, 4, -5},
			out:  []int{2, 4, -5, -5, -9, -66},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var abs []accountbalance.AccountBalance
			for _, amount := range test.in {
				abs = append(abs, accountFn(amount))
			}

			sort.BalanceAmountMagnitude(abs)

			var out []int
			for _, ab := range abs {
				out = append(out, ab.Balance.Amount)
			}

			assert.Equal(t, test.out, out)
		})
	}
}
