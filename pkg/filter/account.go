package filter

import "github.com/glynternet/go-accounting-storage"

type AccountFilter func(storage.Account) bool

func Open() AccountFilter {
	return func(a storage.Account) bool {
		return a.Account.IsOpen()
	}
}

func Filter(as storage.Accounts, f AccountFilter) storage.Accounts {
	var filtered []storage.Account
	for _, a := range as {
		if f(a) {
			filtered = append(filtered, a)
		}
	}
	return storage.Accounts(filtered)
}
