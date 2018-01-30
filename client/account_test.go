package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccountsFromEndpoint(t *testing.T) {
	srv := newJSONResponseTestServer(nil, http.StatusTeapot)
	c := client(srv.URL)
	as, err := c.getAccountsFromEndpoint("")
	assert.Error(t, err)
	assert.Nil(t, as)
}

func newJSONResponseTestServer(encode interface{}, code int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bs, err := json.Marshal(encode)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(code)
		w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
		w.Write(bs)
		return
	}))
}
