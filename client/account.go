package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
	err = json.Unmarshal(bod, as)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response")
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
	err = json.Unmarshal(bod, a)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response")
	}
	return a, err
}

func (c Client) getBodyFromEndpoint(s string) ([]byte, error) {
	res, err := c.getFromEndpoint(s)
	if err != nil {
		return nil, errors.Wrap(err, "getting from endpoint")
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned unexpected code %d (%s)", res.StatusCode, res.Status)
	}
	bod, err := ioutil.ReadAll(res.Body)
	defer func() {
		cErr := res.Body.Close()
		if err == nil {
			err = cErr
		}
	}()
	return bod, errors.Wrap(err, "reading response body")
}

func (c Client) InsertAccount(a account.Account) (*storage.Account, error) {
	return nil, errors.New("not implemented")
}
