package client

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/glynternet/go-accounting/accountingtest"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/common"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/pkg/storage/storagetest"
	"github.com/glynternet/mon/router"
	"github.com/stretchr/testify/assert"
)

const port = 23456

func TestClient_SelectAccounts(t *testing.T) {
	testPort := port + 0

	s := &storagetest.Storage{
		Accounts: &storage.Accounts{
			{
				ID: 51,
				Account: *accountingtest.NewAccount(
					t,
					"test-0",
					accountingtest.NewCurrencyCode(t, "EUR"),
					time.Now().UTC().Truncate(time.Nanosecond),
				),
			},
			{
				ID: 981742,
				Account: *accountingtest.NewAccount(
					t,
					"test-1",
					accountingtest.NewCurrencyCode(t, "GBP"),
					// TODO: Revert this test to not have the UTC() and test when the timezone of the machine running is not UTC.
					time.Now().UTC().Add(time.Hour*123).Truncate(time.Nanosecond),
				),
			},
		},
	}

	r, err := router.New(s)
	assert.NoError(t, err)
	assert.NotNil(t, r)

	srvErr := make(chan error)
	go func() {
		srvErr <- http.ListenAndServe(fmt.Sprintf(":%d", testPort), r)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		selected, err := newTestClient(testPort).SelectAccounts()
		assert.NoError(t, err)
		assert.NotNil(t, selected)
		assert.Equal(t, s.Accounts, selected)
		close(srvErr)
	}()

	common.FatalIfError(t, <-srvErr, "serving")
}

func TestClient_SelectAccount(t *testing.T) {
	testPort := port + 1

	s := &storagetest.Storage{
		Account: &storage.Account{
			ID: 51,
			Account: *accountingtest.NewAccount(
				t,
				"test",
				accountingtest.NewCurrencyCode(t, "EUR"),
				time.Now().UTC().Truncate(time.Nanosecond),
			),
		},
	}

	r, err := router.New(s)
	assert.NoError(t, err)
	assert.NotNil(t, r)

	rErr := make(chan error)
	go func() {
		rErr <- http.ListenAndServe(fmt.Sprintf(":%d", testPort), r)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		selected, err := newTestClient(testPort).SelectAccount(734) // id doesn't matter when mocking
		assert.NoError(t, err)
		assert.NotNil(t, selected)
		assert.Equal(t, s.Account, selected)
		close(rErr)
	}()

	common.FatalIfError(t, <-rErr, "serving")
}

func TestClient_SelectAccountBalances(t *testing.T) {
	testPort := port + 2
	s := &storagetest.Storage{
		Account: &storage.Account{
			ID: 51,
			Account: *accountingtest.NewAccount(
				t,
				"test",
				accountingtest.NewCurrencyCode(t, "EUR"),
				time.Now().Truncate(time.Nanosecond),
			),
		},
		Balances: &storage.Balances{
			storage.Balance{ID: 123},
		},
	}
	r, err := router.New(s)
	common.FatalIfError(t, err, "creating new server")
	assert.NotNil(t, r)

	rErr := make(chan error)
	go func() {
		rErr <- http.ListenAndServe(fmt.Sprintf(":%d", testPort), r)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		selected, err := newTestClient(testPort).SelectAccountBalances(*s.Account) // id doesn't matter when mocking
		assert.NoError(t, err)
		assert.NotNil(t, selected)
		assert.Equal(t, s.Balances, selected)
		close(rErr)
	}()

	common.FatalIfError(t, <-rErr, "serving")
}

func TestClient_InsertAccount(t *testing.T) {
	testPort := port + 3

	account := &storage.Account{
		ID: 51,
		Account: *accountingtest.NewAccount(
			t,
			"test",
			accountingtest.NewCurrencyCode(t, "EUR"),
			time.Date(3000, 0, 0, 0, 0, 0, 0, time.UTC),
		),
	}

	s := &storagetest.Storage{
		Account: account,
	}
	r, err := router.New(s)
	assert.NotNil(t, r)
	common.FatalIfError(t, err, "creating new server")

	srvErr := make(chan error)
	go func() {
		srvErr <- http.ListenAndServe(fmt.Sprintf(":%d", testPort), r)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		inserted, err := newTestClient(testPort).InsertAccount(account.Account)
		assert.NoError(t, err)
		assert.Equal(t, account.ID, inserted.ID)
		assert.Equal(t, account.Account, inserted.Account)
		close(srvErr)
	}()

	common.FatalIfError(t, <-srvErr, "serving")
}

func TestClient_InsertBalance(t *testing.T) {
	testPort := port + 4
	expected := &storage.Balance{ID: 293}

	s := &storagetest.Storage{
		Account: &storage.Account{ID: 51},
		Balance: expected,
	}
	r, err := router.New(s)
	common.FatalIfError(t, err, "creating new server")
	assert.NotNil(t, r)

	rErr := make(chan error)
	go func() {
		rErr <- http.ListenAndServe(fmt.Sprintf(":%d", testPort), r)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		inserted, err := newTestClient(testPort).InsertBalance(storage.Account{}, balance.Balance{})
		assert.NoError(t, err)
		assert.Equal(t, expected, inserted)
		close(rErr)
	}()

	common.FatalIfError(t, <-rErr, "serving")
}

// TODO: newTestClient as closure that increments port everytime it's called
func newTestClient(port int) Client {
	return Client(fmt.Sprintf("http://localhost:%d", port))
}
