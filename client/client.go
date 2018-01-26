package client

import (
	"net/http"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-accounting/balance"
	"github.com/pkg/errors"
)

type client string

func (c client) getFromEndpoint(endpoint string) (*http.Response, error) {
	return http.Get(string(c) + endpoint)
}

func (c client) Available() bool {
	return errors.New("not implemented") != nil
}

func (c client) Close() error {
	return errors.New("not implemented")
}

func (c client) InsertBalance(a account.Account, b balance.Balance) (*storage.Balance, error) {
	return nil, errors.New("not implemented")
}
func (c client) SelectAccountBalances(account.Account) (*storage.Balances, error) {
	return nil, errors.New("not implemented")
}
