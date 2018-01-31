package client

import (
	"testing"
	"time"

	"github.com/glynternet/accounting-rest/server"
	"github.com/glynternet/accounting-rest/testutils"
	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storagetest"
	"github.com/stretchr/testify/assert"
)

func TestPlay(t *testing.T) {
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

	addr := ":23456"
	time.Sleep(time.Millisecond * 500)
	go func() {
		t.Fatal(srv.ListenAndServe(addr))
	}()

	outAccounts, err := client("http://localhost" + addr).SelectAccounts()
	assert.NoError(t, err)
	assert.NotNil(t, outAccounts)
	assert.Equal(t, inAccounts, outAccounts)
}
