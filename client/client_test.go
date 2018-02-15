package client

import (
	"net/http"
	"testing"

	"github.com/glynternet/go-accounting-storage"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// ensure that a Client can be used as a storage.Storage
var _ storage.Storage = Client("")

func Test_getBodyFromEndpoint(t *testing.T) {
	t.Run("get error", func(t *testing.T) {
		c := Client("bloopybloop")
		bod, err := c.getBodyFromEndpoint("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "getting from endpoint")
		assert.Nil(t, bod)
	})

	t.Run("unexpected status", func(t *testing.T) {
		srv := newJSONTestServer(nil, http.StatusTeapot)
		defer srv.Close()
		c := Client(srv.URL)
		as, err := c.getBodyFromEndpoint("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "server returned unexpected code")
		assert.Nil(t, as)
	})
}

type mockMarshal struct {
	err error
}

func (m mockMarshal) MarshalJSON() ([]byte, error) {
	return nil, m.err
}

func Test_postAsJSONToEndpoint(t *testing.T) {
	t.Run("marshal error", func(t *testing.T) {
		c := Client("bloopybloop")
		obj := mockMarshal{
			err: errors.New("can't unmarshal me"),
		}
		res, err := c.postAsJSONToEndpoint("", obj)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "can't unmarshal me")
		}
		assert.Nil(t, res)
	})

	t.Run("post to endpoint error", func(t *testing.T) {
		c := Client("bloopybleep")
		res, err := c.postAsJSONToEndpoint("", nil)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "posting to endpoint")
		}
		assert.Nil(t, res)
	})
}
