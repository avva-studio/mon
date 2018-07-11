package filter

import (
	"time"

	"github.com/glynternet/go-money/currency"
	"github.com/glynternet/mon/pkg/storage"
)

// AccountFilter is a function that will return true if a given storage.Account
// satisfies some certain criteria.
type AccountFilter func(storage.Account) bool

// Existed produces an AccountFilter that can be used to identify if an
// Account existed at a given time
func Existed(t time.Time) AccountFilter {
	return func(a storage.Account) bool {
		return !a.Account.Opened().After(t)
	}
}

// Currencies produces an AccountFilter that will identify a storage.Account if
// it is within a given set of currency codes
func Currencies(cs ...currency.Code) AccountFilter {
	return func(a storage.Account) bool {
		for _, c := range cs {
			if a.Account.CurrencyCode() == c {
				return true
			}
		}
		return false
	}
}

// IDs produces an AccountFilter that will identify a storage.Account if it
// matches on of a given set of IDs
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

// OpenAt produces an AccountFilter that will identify if a storage.Account
// was/is/will-be open at a given time
func OpenAt(t time.Time) AccountFilter {
	return func(a storage.Account) bool {
		return a.Account.OpenAt(t)
	}
}

// Filter returns a set of storage.Accounts that match the given AccountFilter
func Filter(as storage.Accounts, f AccountFilter) storage.Accounts {
	var filtered []storage.Account
	for _, a := range as {
		if f(a) {
			filtered = append(filtered, a)
		}
	}
	return storage.Accounts(filtered)
}
