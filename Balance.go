package main

import (
	"github.com/GlynOwenHanmer/GOHMoneyDB"
	"fmt"
)

type Balance GOHMoneyDB.Balance

// Returns the endpoint location string for updating a balance
func (b Balance) balanceUpdateEndpoint() string {
	return fmt.Sprintf(`/balance/%d/update`, b.Id)
}

