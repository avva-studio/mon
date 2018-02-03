package client

import (
	"testing"
	"time"

	"github.com/glynternet/accounting-rest/server"
	"github.com/glynternet/accounting-rest/testutils"
	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storagetest"
	"github.com/glynternet/go-money/common"
	"github.com/stretchr/testify/assert"
)

const addr = ":23456"

func TestSelectAccounts(t *testing.T) {
	inAccounts := &storage.Accounts{
		{
			Account: accountingtest.NewAccount(
				t,
				"test",
				accountingtest.NewCurrencyCode(t, "EUR"),
				time.Now().Truncate(time.Nanosecond),
			),
		},
	}

	s := &accountingtest.Storage{
		Accounts: inAccounts,
	}

	srv, err := server.New(testutils.NewMockStorageFunc(s, false))
	assert.NoError(t, err)
	assert.NotNil(t, srv)

	srvErr := make(chan error)
	go func() {
		srvErr <- srv.ListenAndServe(addr)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		outAccounts, err := newTestClient().SelectAccounts()
		assert.NoError(t, err)
		assert.NotNil(t, outAccounts)
		assert.Equal(t, inAccounts, outAccounts)
		close(srvErr)
	}()

	common.FatalIfError(t, <-srvErr, "serving")
}

func newTestClient() Client {
	return Client("http://localhost" + addr)
}
