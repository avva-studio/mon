package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/glynternet/accounting-rest/testutils"
	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storagetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_balances(t *testing.T) {
	serveFn := func(s *server, w http.ResponseWriter, r *http.Request) (int, error) {
		return s.balances(1)(w, r)
	}
	nilResponseWriterTest(t, serveFn)
	storageFuncErrorTest(t, serveFn)

	for _, test := range []struct {
		name       string
		code       int
		account    *storage.Account
		accountErr error
		//balance    storage.Balance
		balanceErr error
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
			balanceErr: errors.New("balance error"),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			srv := &server{
				NewStorage: testutils.NewMockStorageFunc(
					&accountingtest.Storage{
						Account:    test.account,
						AccountErr: test.accountErr,
					},
					false,
				),
			}
			code, err := serveFn(srv, rec, nil)
			assert.Equal(t, test.code, code)

			if test.accountErr != nil {
				assert.Equal(t, test.accountErr, errors.Cause(err))
				return
			}
		})
	}
}
