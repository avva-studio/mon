package server

import (
	"net/http"
	"testing"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storagetest"
	"github.com/glynternet/go-accounting/account"
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
			err:  errors.New("selecting handlerSelectAccounts"),
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
	for _, test := range []struct {
		name string
		account.Account
		err  error
		code int
	}{
		{
			name: "zero-values",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := &server{
				storage: &accountingtest.Storage{AccountErr: test.err},
			}
			code, a, err := server.handlerInsertAccount(test.Account)
		})
	}
}
