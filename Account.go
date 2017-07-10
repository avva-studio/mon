package main

import (
	"github.com/GlynOwenHanmer/GOHMoneyDB"
	"fmt"
)

type Account GOHMoneyDB.Account

// Returns the endpoint location string for the balance of the Account
func (account Account) balanceEndpoint() string {
	return fmt.Sprintf(`/account/%d/balance`, account.Id)
}
