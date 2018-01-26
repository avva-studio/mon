package internal

import (
	"fmt"

	"github.com/glynternet/go-accounting-storage"
)

// balance is a wrapper around a DB balance and adds methods for balance endpoints.
type balance storage.Balance

// Returns the endpoint location string for updating a balance
func (b balance) balanceUpdateEndpoint() string {
	return fmt.Sprintf(`/balance/%d/update`, b.ID)
}
