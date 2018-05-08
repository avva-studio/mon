package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storage/storagetest"
	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-accounting/accountingtest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_handlerSelectAccounts(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		expected := errors.New("selecting Accounts")
		server := &server{
			storage: &storagetest.Storage{Err: expected},
		}
		code, as, err := server.handlerSelectAccounts(nil)
		assert.Equal(t, http.StatusServiceUnavailable, code)
		assert.Equal(t, expected, errors.Cause(err))
		assert.Nil(t, as)
	})

	t.Run("success", func(t *testing.T) {
		expected := &storage.Accounts{
			storage.Account{ID: 8767},
		}
		server := &server{
			storage: &storagetest.Storage{
				Accounts: expected,
			},
		}
		code, as, err := server.handlerSelectAccounts(nil)
		assert.Equal(t, http.StatusOK, code)
		assert.NoError(t, err)
		storeAs := as.(*storage.Accounts)
		assert.Equal(t, expected, storeAs)
	})
}

func Test_handlerSelectAccount(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		expected := errors.New("selecting Account")
		server := &server{
			storage: &storagetest.Storage{AccountErr: expected},
		}
		code, a, err := server.handlerSelectAccount(1)
		assert.Equal(t, http.StatusNotFound, code)
		assert.Equal(t, expected, errors.Cause(err))
		assert.Nil(t, a)
	})

	t.Run("success", func(t *testing.T) {
		expected := &storage.Account{
			ID: 456789,
			Account: *accountingtest.NewAccount(t,
				"success account",
				accountingtest.NewCurrencyCode(t, "EUR"),
				time.Date(1000, 0, 0, 0, 0, 0, 0, time.UTC),
			),
		}
		server := &server{
			storage: &storagetest.Storage{Account: expected},
		}
		code, a, err := server.handlerSelectAccount(1)
		assert.Equal(t, http.StatusOK, code)
		assert.NoError(t, err)
		storeA := a.(*storage.Account)
		assert.Equal(t, expected, storeA)
	})
}

func Test_handlerInsertAccount(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		expected := errors.New("inserting account")
		server := &server{
			storage: &storagetest.Storage{AccountErr: expected},
		}
		code, inserted, err := server.handlerInsertAccount(account.Account{})
		assert.Equal(t, expected, errors.Cause(err))
		assert.Nil(t, inserted)
		assert.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("success", func(t *testing.T) {
		expected := &storage.Account{
			ID: 456,
			Account: *accountingtest.NewAccount(t,
				"success account",
				accountingtest.NewCurrencyCode(t, "GBP"),
				time.Date(1000, 1, 0, 0, 0, 0, 0, time.UTC)),
		}
		server := &server{
			storage: &storagetest.Storage{Account: expected},
		}
		code, inserted, err := server.handlerInsertAccount(expected.Account)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, code)
		assert.NotNil(t, inserted)
		assert.Equal(t, expected, inserted)
	})
}
