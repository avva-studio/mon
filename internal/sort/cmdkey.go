package sort

import (
	"github.com/glynternet/mon/internal/accountbalance"
	"github.com/glynternet/mon/pkg/storage"
)

const (
	sortKeyID      = "id"
	sortKeyName    = "name"
	sortKeyBalance = "balance"
)

// AllKeys provides all possible sort keys, agnostic of sort type
func AllKeys() []string {
	ks := []string{""}
	for sortKey, _ := range AccountSorts() {
		ks = append(ks, sortKey)
	}
	for sortKey, _ := range AccountbalanceSorts() {
		ks = append(ks, sortKey)
	}

	return ks
}

// AccountSorts provides a map containing all supported storage.Account sorting
// functions, keyed by the supported keys
func AccountSorts() map[string]func(storage.Accounts) {
	return map[string]func(storage.Accounts){
		sortKeyID:   ID,
		sortKeyName: Name,
	}
}

// AccountbalanceSorts provides a map containing all supported
// accountbalance.AccountBalance sorting functions, keyed by the supported keys
func AccountbalanceSorts() map[string]func([]accountbalance.AccountBalance) {
	return map[string]func([]accountbalance.AccountBalance){
		sortKeyBalance: BalanceAmount,
	}
}
