package main

import (
	"fmt"

	"github.com/GlynOwenHanmer/GOHMoneyDB"
)

type Balance GOHMoneyDB.Balance

// Returns the endpoint location string for updating a balance
func (b Balance) balanceUpdateEndpoint() string {
	return fmt.Sprintf(`/balance/%d/update`, b.Id)
}
