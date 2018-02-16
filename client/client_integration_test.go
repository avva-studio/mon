package client

import (
	"testing"
	"time"

	"fmt"

	"github.com/glynternet/accounting-rest/server"
	"github.com/glynternet/accounting-rest/testutils"
	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storagetest"
	"github.com/glynternet/go-money/common"
	"github.com/stretchr/testify/assert"
)

const port = 23456

func TestClient_SelectAccounts(t *testing.T) {
	testPort := port + 0

	s := &accountingtest.Storage{
		Accounts: &storage.Accounts{
			{
				ID: 51,
				Account: accountingtest.NewAccount(
					t,
					"test-0",
					accountingtest.NewCurrencyCode(t, "EUR"),
					time.Now().Truncate(time.Nanosecond),
				),
			},
			{
				ID: 981742,
				Account: accountingtest.NewAccount(
					t,
					"test-1",
					accountingtest.NewCurrencyCode(t, "GBP"),
					time.Now().Add(time.Hour*123).Truncate(time.Nanosecond),
				),
			},
		},
	}

	srv, err := server.New(testutils.NewMockStorageFunc(s, false))
	assert.NoError(t, err)
	assert.NotNil(t, srv)

	srvErr := make(chan error)
	go func() {
		srvErr <- srv.ListenAndServe(fmt.Sprintf(":%d", testPort))
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

	s := &accountingtest.Storage{
		Account: &storage.Account{
			ID: 51,
			Account: accountingtest.NewAccount(
				t,
				"test",
				accountingtest.NewCurrencyCode(t, "EUR"),
				time.Now().Truncate(time.Nanosecond),
			),
		},
	}

	srv, err := server.New(testutils.NewMockStorageFunc(s, false))
	assert.NoError(t, err)
	assert.NotNil(t, srv)

	srvErr := make(chan error)
	go func() {
		srvErr <- srv.ListenAndServe(fmt.Sprintf(":%d", testPort))
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		selected, err := newTestClient(testPort).SelectAccount(734) // id doesn't matter when mocking
		assert.NoError(t, err)
		assert.NotNil(t, selected)
		assert.Equal(t, s.Account, selected)
		close(srvErr)
	}()

	common.FatalIfError(t, <-srvErr, "serving")
}

func TestClient_SelectAccountBalances(t *testing.T) {
	testPort := port + 2
	s := &accountingtest.Storage{
		Account: &storage.Account{
			ID: 51,
			Account: accountingtest.NewAccount(
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
	srv, err := server.New(testutils.NewMockStorageFunc(s, false))
	assert.NoError(t, err)
	assert.NotNil(t, srv)

	srvErr := make(chan error)
	go func() {
		srvErr <- srv.ListenAndServe(fmt.Sprintf(":%d", testPort))
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		selected, err := newTestClient(testPort).SelectAccountBalances(*s.Account) // id doesn't matter when mocking
		assert.NoError(t, err)
		assert.NotNil(t, selected)
		assert.Equal(t, s.Balances, selected)
		close(srvErr)
	}()

	common.FatalIfError(t, <-srvErr, "serving")
}

// TODO: Do these tests actually want to be split up?
// What are the best practices here?
func TestClient_InsertAccount2(t *testing.T) {
	testPort := port + 3

	account := &storage.Account{
		ID: 51,
		Account: accountingtest.NewAccount(
			t,
			"test",
			accountingtest.NewCurrencyCode(t, "EUR"),
			time.Now().Truncate(time.Nanosecond),
		),
	}

	s := &accountingtest.Storage{
		Account: account,
	}
	srv, err := server.New(testutils.NewMockStorageFunc(s, false))
	common.FatalIfError(t, err, "creating mock storage")
	assert.NotNil(t, srv)

	srvErr := make(chan error)
	go func() {
		srvErr <- srv.ListenAndServe(fmt.Sprintf(":%d", testPort))
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		inserted, err := newTestClient(testPort).InsertAccount(account.Account)
		assert.NoError(t, err)
		assert.NotNil(t, inserted)
		//assert.Equal(t, account.ID, inserted.ID)
		//assert.Equal(t, account.Account, inserted.Account)
		close(srvErr)
	}()

	common.FatalIfError(t, <-srvErr, "serving")
}

// TODO: newTestClient as closure that increments port everytime it's called
func newTestClient(port int) Client {
	return Client(fmt.Sprintf("http://localhost:%d", port))
}
