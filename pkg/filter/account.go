package filter

import (
	"time"

	"github.com/glynternet/mon/pkg/storage"
)

type AccountFilter func(storage.Account) bool

func Existed(t time.Time) AccountFilter {
	return func(a storage.Account) bool {
		return !a.Account.Opened().After(t)
	}
}

func IDs(ids []uint) AccountFilter {
	return func(a storage.Account) bool {
		for _, id := range ids {
			if a.ID == id {
				return true
			}
		}
		return false
	}
}

func OpenAt(t time.Time) AccountFilter {
	return func(a storage.Account) bool {
		return a.Account.OpenAt(t)
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
