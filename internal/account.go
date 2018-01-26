package internal

import (
	"fmt"

	"github.com/glynternet/go-accounting-storage"
)

// account is a wrapper around a DB account to add methods for certain endpoints.
type account storage.Account

// Returns the endpoint location string for the balance of the account
func (a account) balanceEndpoint() string {
	return a.generateEndpoint(`balance`)
}

// Returns the endpoint location string for getting the balances of an account
func (a account) balancesEndpoint() string {
	return a.generateEndpoint(`balances`)
}

// Returns the endpoint location string for updating an account
func (a account) updateEndpoint() string {
	return a.generateEndpoint(`update`)
}

// Returns the endpoint location string for deleting the account
func (a account) deleteEndpoint() string {
	return a.generateEndpoint(`delete`)
}

func (a account) generateEndpoint(endpoint string) string {
	return fmt.Sprintf(`/account/%d/%s`, a.ID, endpoint)
}
