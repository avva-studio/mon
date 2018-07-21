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

// AllSortKeys provides all possible sort keys, agnostic of sort type
func AllSortKeys() Keys {
	var ks []string
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

// Keys is a set of string that can be used as keys for sorting
type Keys []string

// Contains identifies if a given key/string is in a set of Keys
func (ks Keys) Contains(key string) bool {
	for _, valid := range ks {
		if valid == key {
			return true
		}
	}
	return false
}
