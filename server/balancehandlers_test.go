package server

import (
	"net/http"
	"testing"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storagetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_balances(t *testing.T) {
	t.Run("SelectAccount error", func(t *testing.T) {
		expected := errors.New("account error")
		srv := &server{
			storage: &accountingtest.Storage{
				AccountErr: expected,
			},
		}
		code, bs, err := srv.balances(1) // any ID can be used because of the stub
		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, expected, errors.Cause(err))
		assert.Nil(t, bs)
	})

	t.Run("SelectBalance error", func(t *testing.T) {
		account := &storage.Account{ID: 51}
		expected := errors.New("balances error")
		srv := &server{
			storage: &accountingtest.Storage{
				Account:     account,
				BalancesErr: expected,
			},
		}
		code, bs, err := srv.balances(1) // any ID can be used because of the stub
		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, expected, errors.Cause(err))
		assert.Nil(t, bs)
	})

	t.Run("all ok", func(t *testing.T) {
		account := &storage.Account{ID: 51}
		srv := &server{
			storage: &accountingtest.Storage{
				Account: account,
			},
		}
		code, bs, err := srv.balances(1) // any ID can be used because of the stub
		assert.Equal(t, http.StatusOK, code)
		assert.NoError(t, err)
		assert.IsType(t, new(storage.Balances), bs)
	})
}
