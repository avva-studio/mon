package server

import (
	"net/http"
	"testing"

	"time"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storagetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_handlerSelectAccounts(t *testing.T) {
	for _, test := range []struct {
		name string
		code int
		err  error
	}{
		{
			name: "error",
			code: http.StatusServiceUnavailable,
			err:  errors.New("selecting handlerSelectAccounts"),
		},
		{
			name: "success",
			code: http.StatusOK,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := &server{
				storage: &accountingtest.Storage{Err: test.err},
			}
			code, as, err := server.handlerSelectAccounts(nil)
			assert.Equal(t, test.code, code)

			if test.err != nil {
				assert.Equal(t, test.err, errors.Cause(err))
				return
			}
			assert.NoError(t, err)
			_, ok := as.(*storage.Accounts)
			assert.True(t, ok)
		})
	}
}

func Test_handlerSelectAccount(t *testing.T) {
	for _, test := range []struct {
		name string
		code int
		err  error
	}{
		{
			name: "error",
			code: http.StatusNotFound,
			err:  errors.New("selecting Account"),
		},
		{
			name: "success",
			code: http.StatusOK,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := &server{
				storage: &accountingtest.Storage{AccountErr: test.err},
			}
			code, a, err := server.handlerSelectAccount(1)
			assert.Equal(t, test.code, code)

			if test.err != nil {
				assert.Equal(t, test.err, errors.Cause(err))
				return
			}

			assert.NoError(t, err)
			_, ok := a.(*storage.Account)
			assert.True(t, ok)
		})
	}
}

func Test_handlerInsertAccount(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		expected := errors.New("inserting account")
		server := &server{
			storage: &accountingtest.Storage{AccountErr: expected},
		}
		code, inserted, err := server.handlerInsertAccount(nil)
		assert.Equal(t, expected, errors.Cause(err))
		assert.Nil(t, inserted)
		assert.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("success", func(t *testing.T) {
		expected := &storage.Account{
			ID: 456,
			Account: accountingtest.NewAccount(t,
				"success account",
				accountingtest.NewCurrencyCode(t, "GBP"),
				time.Date(1000, 1, 0, 0, 0, 0, 0, time.UTC)),
		}
		server := &server{
			storage: &accountingtest.Storage{Account: expected},
		}
		code, inserted, err := server.handlerInsertAccount(expected.Account)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, code)
		assert.NotNil(t, inserted)
		assert.Equal(t, expected, inserted)
	})

}
