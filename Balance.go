package main

import (
	"fmt"

	"github.com/glynternet/GOHMoneyDB"
)

// Balance is a wrapper around a DB Balance and adds methods for Balance endpoints.
type Balance GOHMoneyDB.Balance

// Returns the endpoint location string for updating a balance
func (b Balance) balanceUpdateEndpoint() string {
	return fmt.Sprintf(`/balance/%d/update`, b.ID)
}
