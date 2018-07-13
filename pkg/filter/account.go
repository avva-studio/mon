package filter

import (
	"time"

	"github.com/glynternet/go-money/currency"
	"github.com/glynternet/mon/pkg/storage"
)

// AccountCondition is a function that will return true if a given storage.Account
// satisfies some certain condition.
type AccountCondition func(storage.Account) bool

// Existed produces an AccountCondition that can be used to identify if an
// Account existed/exists/will-exist at a given time.
// Existed will identify that an Account existed if its open date matches or
// was before the given time
func Existed(t time.Time) AccountCondition {
	return func(a storage.Account) bool {
		return !a.Account.Opened().After(t)
	}
}

// OpenAt produces an AccountCondition that will identify if a storage.Account
// was/is/will-be open at a given time
func OpenAt(t time.Time) AccountCondition {
	return func(a storage.Account) bool {
		return a.Account.OpenAt(t)
	}
}

// Currency produces an AccountCondition that will identify a storage.Account if
// it has a given currency.Code
func Currency(c currency.Code) AccountCondition {
	return func(a storage.Account) bool {
		return a.Account.CurrencyCode() == c
	}
}

// ID produces an AccountCondition that will identify a storage.Account if it
// matches an ID
func ID(id uint) AccountCondition {
	return func(a storage.Account) bool {
		return a.ID == id
	}
}

// AccountConditions is a set of AccountCondition
type AccountConditions []AccountCondition

// Or identifies when an account satisfies one of more constraints of an
// AccountConditions
func (afs AccountConditions) Or(a storage.Account) bool {
	for _, af := range afs {
		if af(a) {
			return true
		}
	}
	return false
}

// Filter returns a set of storage.Accounts that match the given AccountCondition
func Filter(as storage.Accounts, f AccountCondition) storage.Accounts {
	var filtered []storage.Account
	for _, a := range as {
		if f(a) {
			filtered = append(filtered, a)
		}
	}
	return storage.Accounts(filtered)
}
