package client

import (
	"encoding/json"

	"fmt"

	"github.com/glynternet/accounting-rest/server"
	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting/balance"
	"github.com/pkg/errors"
)

func (c Client) SelectAccountBalances(a storage.Account) (*storage.Balances, error) {
	return c.getBalancesFromEndpoint(fmt.Sprintf(server.EndpointFmtAccountBalances, a.ID))
}

func (c Client) getBalancesFromEndpoint(e string) (*storage.Balances, error) {
	bod, err := c.getBodyFromEndpoint(e)
	if err != nil {
		return nil, errors.Wrap(err, "getting body from endpoint")
	}
	bs := new(storage.Balances)
	err = errors.Wrap(json.Unmarshal(bod, bs), "unmarshalling response")
	if err != nil {
		bs = nil
	}
	return bs, err
}

func (c Client) InsertBalance(a storage.Account, b balance.Balance) (*storage.Balance, error) {
	return nil, errors.New("not implemented")
}
