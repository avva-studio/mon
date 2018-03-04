package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountsFromEndpoint(t *testing.T) {
	t.Run("get body error", func(t *testing.T) {
		c := Client("bloopybloop")
		as, err := c.getAccountsFromEndpoint("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "getting from endpoint")
		assert.Nil(t, as)
	})

	t.Run("unmarshallable response error", func(t *testing.T) {
		srv := newJSONTestServer(
			struct{ NonAccount string }{NonAccount: "bloop"},
			http.StatusOK,
		)
		defer srv.Close()
		c := Client(srv.URL)
		as, err := c.getAccountsFromEndpoint("")
		if assert.Error(t, err) {
			assert.IsType(t, &json.UnmarshalTypeError{}, errors.Cause(err))
		}
		assert.Nil(t, as)
	})
}

func TestGetAccountFromEndpoint(t *testing.T) {
	t.Run("get body error", func(t *testing.T) {
		c := Client("bloopybleep")
		a, err := c.getAccountFromEndpoint("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "getting from endpoint")
		assert.Nil(t, a)
	})

	t.Run("unmarshallable response", func(t *testing.T) {
		srv := newJSONTestServer(
			struct{ NonAccount string }{NonAccount: "bloop"},
			http.StatusOK,
		)
		defer srv.Close()
		c := Client(srv.URL)
		as, err := c.getAccountFromEndpoint("")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "json unmarshalling into account")
		}
		assert.Nil(t, as)
	})
}

func TestPostAccountToEndpoint(t *testing.T) {
	t.Run("post as json error", func(t *testing.T) {
		bod, err := Client("BLOOOOP").postAccountToEndpoint("", nil)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "posting as JSON")
		}
		assert.Nil(t, bod)
	})
	//t.Run("processRequestForBody err", func(t *testing.T) {
	//	_ := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		 TODO: make a processRequestForBody error test
	//}))
	//})
}

//func TestClient_InsertAccount(t *testing.T) {
//	t.Run("postAccountToEndpoint error", func(t *testing.T) {
//		a := accountingtest.NewAccount(t,
//			"any",
//			accountingtest.NewCurrencyCode(t, "SKH"),
//			time.Now())
//		c := Client("WOOOOOOH")
//		storageA, err := c.InsertAccount(a)
//		if assert.Error(t, err) {
//			assert.Contains(t, err.Error(), "posting account to endpoint")
//		}
//		assert.Nil(t, storageA)
//	})
//
//	t.Run("unmarshallable response", func(t *testing.T) {
//		srv := newJSONTestServer(
//			struct{ NonAccount string }{NonAccount: "bloop"},
//			http.StatusOK,
//		)
//		defer srv.Close()
//		c := Client(srv.URL)
//		as, err := c.InsertAccount(nil)
//		if assert.Error(t, err) {
//			assert.Contains(t, err.Error(), "json unmarshalling response")
//		}
//		assert.Nil(t, as)
//	})
//}

func newJSONTestServer(encode interface{}, code int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bs, err := json.Marshal(encode)
		if err != nil {
			panic(fmt.Sprintf("error marshalling to json: %v", err))
		}
		w.WriteHeader(code)
		w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
		_, err = w.Write(bs)
		if err != nil {
			panic(fmt.Sprintf("error writing to ResponseWriter: %v", err))
		}
	}))
}
