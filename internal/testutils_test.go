package internal

import (
	"github.com/glynternet/go-accounting-storage"
	aaccount "github.com/glynternet/go-accounting/account"
	abalance "github.com/glynternet/go-accounting/balance"
)

type mockStorage struct {
	available bool
	err       error
	*storage.Account
	*storage.Accounts
	*storage.Balance
	*storage.Balances
}

func (s mockStorage) Available() bool { return s.available }
func (s mockStorage) Close() error    { return s.err }

func (s mockStorage) InsertAccount(a aaccount.Account) (*storage.Account, error) {
	return s.Account, s.err
}
func (s mockStorage) SelectAccount(u uint) (*storage.Account, error) { return s.Account, s.err }
func (s mockStorage) SelectAccounts() (*storage.Accounts, error)     { return s.Accounts, s.err }

func (s mockStorage) InsertBalance(a storage.Account, b abalance.Balance) (*storage.Balance, error) {
	return s.Balance, s.err
}
func (s mockStorage) SelectAccountBalances(storage.Account) (*storage.Balances, error) {
	return s.Balances, s.err
}

func (s mockStorage) storageFunc() (storage.Storage, error) {
	return s, nil
}
