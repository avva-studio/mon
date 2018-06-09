package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Client is a client to retrieve accounting items over http using REST
type Client string

func defaultHTTPClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}

func (c Client) getFromEndpoint(endpoint string) (*http.Response, error) {
	return http.Get(string(c) + endpoint)
}

func (c Client) postToEndpoint(endpoint string, contentType string, body io.Reader) (*http.Response, error) {
	return http.Post(string(c)+endpoint, contentType, body)
}

func (c Client) deleteToEndpoint(endpoint string) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodDelete, string(c)+endpoint, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating new request")
	}
	return defaultHTTPClient().Do(r)
}

func (c Client) Available() bool {
	// TODO: Deprecate Available in favour of something that returns more information
	_, err := c.SelectAccounts()
	return err == nil
}

func (c Client) Close() error {
	return nil
}

func (c Client) getBodyFromEndpoint(e string) ([]byte, error) {
	res, err := c.getFromEndpoint(e)
	if err != nil {
		return nil, errors.Wrap(err, "getting from endpoint")
	}
	return processResponseForBody(res)
}

func processResponseForBody(r *http.Response) ([]byte, error) {
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned unexpected code %d (%s)", r.StatusCode, r.Status)
	}
	bod, err := ioutil.ReadAll(r.Body)

	defer func() {
		// TODO: this handler only needs to take a []byte which would mean we can handle closing the body elsewhere
		cErr := r.Body.Close()
		if cErr != nil {
			log.Print(errors.Wrap(err, "closing response body"))
		}
	}()

	return bod, errors.Wrap(err, "reading response body")
}

func (c Client) postAsJSONToEndpoint(e string, thing interface{}) (*http.Response, error) {
	bs, err := json.Marshal(thing)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling json")
	}
	res, err := c.postToEndpoint(e, `application/json; charset=UTF-8`, bytes.NewReader(bs))
	return res, errors.Wrap(err, "posting to endpoint")
}
