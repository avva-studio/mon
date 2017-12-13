package internal

import (
	"fmt"

	"github.com/glynternet/go-accounting-storage"
)

// Account is a wrapper around a DB Account to add methods for certain endpoints.
type Account storage.Account

// Returns the endpoint location string for the balance of the Account
func (a Account) balanceEndpoint() string {
	return a.generateEndpoint(`balance`)
}

// Returns the endpoint location string for getting the balances of an Account
func (a Account) balancesEndpoint() string {
	return a.generateEndpoint(`balances`)
}

// Returns the endpoint location string for updating an Account
func (a Account) updateEndpoint() string {
	return a.generateEndpoint(`update`)
}

// Returns the endpoint location string for deleting the Account
func (a Account) deleteEndpoint() string {
	return a.generateEndpoint(`delete`)
}

func (a Account) generateEndpoint(endpoint string) string {
	return fmt.Sprintf(`/account/%d/%s`, a.ID, endpoint)
}
