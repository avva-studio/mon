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
// Account existed/exists/will-exist at a given time
func Existed(t time.Time) AccountFilter {
	return func(a storage.Account) bool {
		return !a.Account.Opened().After(t)
	}
}

// OpenAt produces an AccountFilter that will identify if a storage.Account
// was/is/will-be open at a given time
func OpenAt(t time.Time) AccountFilter {
	return func(a storage.Account) bool {
		return a.Account.OpenAt(t)
	}
}

// Currency produces an AccountFilter that will identify a storage.Account if
// it has a given currency.Code
func Currency(c currency.Code) AccountFilter {
	return func(a storage.Account) bool {
		return a.Account.CurrencyCode() == c
	}
}

// ID produces an AccountFilter that will identify a storage.Account if it
// matches an ID
func ID(id uint) AccountFilter {
	return func(a storage.Account) bool {
		return a.ID == id
	}
}

// AccountFilters is a set of AccountFilter
type AccountFilters []AccountFilter

// Or identifies when an account satisfies one of more constaints of an
// AccountFilters
func (afs AccountFilters) Or(a storage.Account) bool {
	for _, af := range afs {
		if af(a) {
			return true
		}
	}
	return false
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
