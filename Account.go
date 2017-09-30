package main

import (
	"fmt"

	"github.com/GlynOwenHanmer/GOHMoneyDB"
)

// Account is a wrapper around a DB Account to add methods for certain endpoints.
type Account GOHMoneyDB.Account

// Returns the endpoint location string for the balance of the Account
func (a Account) balanceEndpoint() string {
	return fmt.Sprintf(`/account/%d/balance`, a.ID)
}

// Returns the endpoint location string for updating an Account
func (a Account) updateEndpoint() string {
	return fmt.Sprintf(`/account/%d/update`, a.ID)
}

// Returns the endpoint location string for deleting the Account
func (a Account) deleteEndpoint() string {
	return fmt.Sprintf(`/account/%d/delete`, a.ID)
}
