package sort

import (
	"sort"

	"github.com/glynternet/mon/internal/accountbalance"
)

// BalanceAmount sorts a slice of accountbalance.AccountBalance by the amount
// of the Balance in ascending order.
// BalanceAmount cannot guarantee any specific order within a subsection of
// the slice when multiple AccountBalance have the same amount.
func BalanceAmount(abs []accountbalance.AccountBalance) {
	sort.Slice(abs, func(i, j int) bool {
		return (abs)[i].Balance.Amount < (abs)[j].Balance.Amount
	})
}

func BalanceAmountMagnitude(abs []accountbalance.AccountBalance) {
	sort.Slice(abs, func(i, j int) bool {
		absI := absolute((abs)[i].Balance.Amount)
		absJ := absolute((abs)[j].Balance.Amount)
		return absI < absJ
	})
}

func absolute(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
