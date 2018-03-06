package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glynternet/go-accounting/balance"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetBalancesFromEndpoint(t *testing.T) {
	t.Run("get body error", func(t *testing.T) {
		c := Client("bloopybloop")
		as, err := c.getBalancesFromEndpoint("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "getting from endpoint")
		assert.Nil(t, as)
	})

	t.Run("unmarshallable response", func(t *testing.T) {
		srv := newJSONTestServer(
			struct{ NonAccount string }{NonAccount: "bloop"},
			http.StatusOK,
		)
		defer srv.Close()
		c := Client(srv.URL)
		bs, err := c.getBalancesFromEndpoint("")
		if assert.Error(t, err) {
			assert.IsType(t, &json.UnmarshalTypeError{}, errors.Cause(err))
		}
		assert.Nil(t, bs)
	})
}

func TestClient_postBalanceToEndpoint(t *testing.T) {
	t.Run("post as json error", func(t *testing.T) {
		bod, err := Client("BLOOOOP").postBalanceToEndpoint("", balance.Balance{})
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "posting as JSON")
		}
		assert.Nil(t, bod)
	})

	t.Run("processResponseForBody err", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		bod, err := Client(srv.URL).postBalanceToEndpoint("", balance.Balance{})
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "server returned unexpected code ")
		}
		assert.Empty(t, bod)
	})
}
