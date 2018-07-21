package accountbalance

import (
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/mon/pkg/storage"
)

// AccountBalance represents the state of a storage.Account at a given moment,
// the moment in time being determined by the time of the balance.Balance.
type AccountBalance struct {
	storage.Account
	balance.Balance
}
