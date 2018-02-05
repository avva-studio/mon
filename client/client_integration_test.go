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

func newTestClient(port int) Client {
	return Client(fmt.Sprintf("http://localhost:%d", port))
}
