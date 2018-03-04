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

func Test_balances(t *testing.T) {
	for _, test := range []struct {
		name        string
		code        int
		account     *storage.Account
		accountErr  error
		balancesErr error
	}{
		{
			name:       "SelectAccount error",
			code:       http.StatusBadRequest,
			accountErr: errors.New("account error"),
		},
		{
			name: "SelectBalance error",
			code: http.StatusBadRequest,
			account: &storage.Account{
				ID: 51,
				Account: accountingtest.NewAccount(
					t,
					"test",
					accountingtest.NewCurrencyCode(t, "EUR"),
					time.Now().Truncate(time.Nanosecond),
				),
			},
			balancesErr: errors.New("balances error"),
		},
		{
			name: "all ok",
			code: http.StatusOK,
			account: &storage.Account{
				ID: 51,
				Account: accountingtest.NewAccount(
					t,
					"test",
					accountingtest.NewCurrencyCode(t, "EUR"),
					time.Now().Truncate(time.Nanosecond),
				),
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			srv := &server{
				storage: &accountingtest.Storage{
					Account:     test.account,
					AccountErr:  test.accountErr,
					BalancesErr: test.balancesErr,
				},
			}
			code, bs, err := srv.balances(1)
			assert.Equal(t, test.code, code)

			if test.accountErr != nil {
				assert.Equal(t, test.accountErr, errors.Cause(err))
				return
			}
			if test.balancesErr != nil {
				assert.Equal(t, test.balancesErr, errors.Cause(err))
				return
			}
			assert.IsType(t, new(storage.Balances), bs)
		})
	}
}
