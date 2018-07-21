package sort

import (
	"sort"

	"github.com/glynternet/mon/pkg/storage"
)

// ID sorts a storage.Accounts by ID in ascending order.
// ID cannot guarantee any specific order within a subsection of
// storage.Accounts when multiple accounts have the same ID.
func ID(as storage.Accounts) {
	sort.Slice(as, func(i, j int) bool {
		return (as)[i].ID < (as)[j].ID
	})
}

// Name sorts a storage.Accounts by Name in ascending order.
// Name cannot guarantee any specific order within a subsection of
// storage.Accounts when multiple accounts have the same Name.
func Name(as storage.Accounts) {
	sort.Slice(as, func(i, j int) bool {
		return (as)[i].Account.Name() < (as)[j].Account.Name()
	})
}
