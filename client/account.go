package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/glynternet/accounting-rest/internal"
	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting/account"
	"github.com/pkg/errors"
)

func (c client) Accounts() (*storage.Accounts, error) {
	return c.getAccountsFromEndpoint(internal.EndpointAccounts)
}

func (c client) getAccountsFromEndpoint(e string) (*storage.Accounts, error) {
	res, err := c.getFromEndpoint(internal.EndpointAccounts)
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
	defer res.Body.Close()
	as := new(storage.Accounts)
	err = json.Unmarshal(bod, as)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response")
	}
	return as, err
}

func (c client) InsertAccount(a account.Account) (*storage.Account, error) {
	return nil, errors.New("not implemented")
}

func (c client) SelectAccount(u uint) (*storage.Account, error) {
	return nil, errors.New("not implemented")
}

func (c client) SelectAccounts() (*storage.Accounts, error) {
	return nil, errors.New("not implemented")
}
