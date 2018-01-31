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
	res, err := c.getFromEndpoint(server.EndpointAccounts)
	if err != nil {
		return nil, errors.Wrap(err, "getting from endpoint")
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned code %d (%s)", res.StatusCode, res.Status)
	}
	bod, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body")
	}
	defer func() {
		cErr := res.Body.Close()
		if err == nil {
			err = cErr
		}
	}()
	as := new(storage.Accounts)
	err = json.Unmarshal(bod, as)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response")
	}
	return as, err
}

func (c Client) SelectAccount(u uint) (*storage.Account, error) {
	return nil, errors.New("not implemented")
}

func (c Client) InsertAccount(a account.Account) (*storage.Account, error) {
	return nil, errors.New("not implemented")
}
