package main

import (
	"fmt"

	"github.com/glynternet/go-accounting-storage"
)

// Balance is a wrapper around a DB Balance and adds methods for Balance endpoints.
type Balance storage.Balance

// Returns the endpoint location string for updating a balance
func (b Balance) balanceUpdateEndpoint() string {
	return fmt.Sprintf(`/balance/%d/update`, b.ID)
}
