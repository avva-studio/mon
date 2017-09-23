package main

import (
	"fmt"

	"github.com/GlynOwenHanmer/GOHMoneyDB"
)

type Account GOHMoneyDB.Account

// Returns the endpoint location string for the balance of the Account
func (a Account) balanceEndpoint() string {
	return fmt.Sprintf(`/account/%d/balance`, a.Id)
}

// Returns the endpoint location string for updating an Account
func (a Account) updateEndpoint() string {
	return fmt.Sprintf(`/account/%d/update`, a.Id)
}
