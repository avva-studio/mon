package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/server"
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
	as := &storage.Accounts{}
	err = errors.Wrapf(json.Unmarshal(bod, as), "unmarshalling response body: %s", string(bod))
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
	return unmarshalJSONToAccount(bod)
}

func (c Client) InsertAccount(a account.Account) (*storage.Account, error) {
	bs, err := c.postAccountToEndpoint(server.EndpointAccountInsert, a)
	if err != nil {
		return nil, errors.Wrapf(err, "posting account to endpoint %s", server.EndpointAccountInsert)
	}
	return unmarshalJSONToAccount(bs)
}

func (c Client) UpdateAccount(account *storage.Account, updates *account.Account) (*storage.Account, error) {
	endpoint := fmt.Sprintf(server.EndpointFmtAccountUpdate, account.ID)
	bs, err := c.postAccountToEndpoint(endpoint, *updates)
	if err != nil {
		return nil, errors.Wrapf(err, "posting account to endpoint %s", endpoint)
	}
	return unmarshalJSONToAccount(bs)
}

func (c Client) DeleteAccount(id uint) error {
	endpoint := fmt.Sprintf(server.EndpointFmtAccount, id)
	r, err := c.deleteToEndpoint(endpoint)
	if err != nil {
		return errors.Wrapf(err, "deleting account to endpoint %s", endpoint)
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d (%s)", r.StatusCode, http.StatusText(r.StatusCode))
	}
	return nil
}

func (c Client) postAccountToEndpoint(e string, a account.Account) ([]byte, error) {
	res, err := c.postAsJSONToEndpoint(e, a)
	if err != nil {
		return nil, errors.Wrap(err, "posting as JSON")
	}
	return processResponseForBody(res)
}

func unmarshalJSONToAccount(data []byte) (*storage.Account, error) {
	if len(data) == 0 {
		return nil, errors.New("no data provided")
	}
	if string(data) == "null" {
		return nil, errors.New(`data was "null"`)
	}
	a := &storage.Account{}
	err := errors.Wrapf(json.Unmarshal(data, a), "json unmarshalling into account. bytes as string: %s", data)
	if err != nil {
		a = nil
	}
	return a, err
}
