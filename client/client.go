package client

import (
	"net/http"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting/balance"
	"github.com/pkg/errors"
)

// Client is a client to retrieve accounting items over http using REST
type Client string

func (c Client) getFromEndpoint(endpoint string) (*http.Response, error) {
	return http.Get(string(c) + endpoint)
}

func (c Client) Available() bool {
	return errors.New("not implemented") != nil
}

func (c Client) Close() error {
	return errors.New("not implemented")
}

func (c Client) InsertBalance(a storage.Account, b balance.Balance) (*storage.Balance, error) {
	return nil, errors.New("not implemented")
}
func (c Client) SelectAccountBalances(storage.Account) (*storage.Balances, error) {
	return nil, errors.New("not implemented")
}
