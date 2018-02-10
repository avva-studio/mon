package client

import (
	"net/http"

	"fmt"
	"io/ioutil"

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
