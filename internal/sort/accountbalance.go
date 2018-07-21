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
