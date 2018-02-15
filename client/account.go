package client

import (
	"encoding/json"
	"fmt"

	"github.com/glynternet/accounting-rest/server"
	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting/account"
	"github.com/pkg/errors"
)

func (c Client) SelectAccounts() (*storage.Accounts, error) {
	return c.getAccountsFromEndpoint(server.EndpointAccounts)
}

func (c Client) getAccountsFromEndpoint(e string) (*storage.Accounts, error) {
	bod, err := c.getBodyFromEndpoint(e)
	if err != nil {
		return nil, errors.Wrap(err, "getting body from endpoint")
	}
	as := new(storage.Accounts)
	err = errors.Wrap(json.Unmarshal(bod, as), "unmarshalling response")
	if err != nil {
		as = nil
	}
	return as, err
}

func (c Client) SelectAccount(u uint) (*storage.Account, error) {
	return c.getAccountFromEndpoint(fmt.Sprintf(server.EndpointFmtAccount, u))
}

func (c Client) getAccountFromEndpoint(e string) (*storage.Account, error) {
	bod, err := c.getBodyFromEndpoint(e)
	if err != nil {
		return nil, errors.Wrap(err, "getting body from endpoint")
	}
	a := new(storage.Account)
	err = errors.Wrap(json.Unmarshal(bod, a), "unmarshalling response")
	if err != nil {
		a = nil
	}
	return a, err
}

func (c Client) InsertAccount(a account.Account) (*storage.Account, error) {
	_, err := c.postAccountToEndpoint(server.EndpointAccountInsert, a)
	if err != nil {
		return nil, errors.Wrapf(err, "posting account to endpoint %f", server.EndpointAccountInsert)
	}
	return nil, nil
}

func (c Client) postAccountToEndpoint(e string, a account.Account) ([]byte, error) {
	res, err := c.postAsJSONToEndpoint(e, a)
	if err != nil {
		return nil, errors.Wrap(err, "posting as JSON")
	}
	bod, err := getResponseBody(res)
	return bod, errors.Wrap(err, "getting response body")
}
