package storage

import (
	"github.com/glynternet/go-accounting/balance"
)

// Balance holds logic for an Account item that is held within a go-money database.
type Balance struct {
	balance.Balance
	ID uint
}

// Equal returns true if two Balance items are logically identical
func (b Balance) Equal(ob Balance) bool {
	if b.ID != ob.ID || !b.Balance.Equal(ob.Balance) {
		return false
	}
	return true
}

// Balances holds multiple Balance items
type Balances []Balance

// InnerBalances returns the balance.Balances contained within a set of Balances
func (bs Balances) InnerBalances() balance.Balances {
	var bbs balance.Balances
	for _, b := range bs {
		bbs = append(bbs, b.Balance)
	}
	return bbs
}
