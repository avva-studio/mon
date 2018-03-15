package functional

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storage/postgres2"
	"github.com/glynternet/go-accounting-storagetest"
	"github.com/glynternet/go-money/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	// viper keys
	keyDBHost    = "db-host"
	keyDBUser    = "db-user"
	keyDBName    = "db-name"
	keyDBSSLMode = "db-sslmode"
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
}

func TestCreateStorage(t *testing.T) {
	const retries = 5

	var err error
	var i int
	for i = 1; i <= retries; i++ {
		err = postgres2.CreateStorage(
			viper.GetString(keyDBHost),
			viper.GetString(keyDBUser),
			viper.GetString(keyDBName),
			viper.GetString(keyDBSSLMode),
		)
		if err == nil {
			break
		}
		log.Printf("Attempt: %d, err: %v\n", i, err)
		time.Sleep(time.Second)
	}
	assert.NoError(t, err)
	t.Logf("Attempts: %d", i)
}

func TestInsertingAndRetrievingTwoAccounts(t *testing.T) {
	db := createStorage(t)

	as, err := db.SelectAccounts()
	common.FatalIfError(t, err, "selecting accounts")

	if !assert.Len(t, *as, 0) {
		t.FailNow()
	}

	a := accountingtest.NewAccount(t, "A", accountingtest.NewCurrencyCode(t, "BTC"), time.Now())
	insertedA, err := db.InsertAccount(a)
	common.FatalIfError(t, err, "inserting account")

	as, err = db.SelectAccounts()
	common.FatalIfError(t, err, "selecting accounts after inserting one")

	if !assert.Len(t, *as, 1) {
		t.FailNow()
	}
	retrievedA := (*as)[0]
	equal, err := insertedA.Equal(retrievedA)
	common.FatalIfError(t, err, "equaling inserted and retrieved")
	if !assert.True(t, equal) {
		t.FailNow()
	}

	b := accountingtest.NewAccount(t, "B", accountingtest.NewCurrencyCode(t, "EUR"), time.Now().Add(-1*time.Hour))

	insertedB, err := db.InsertAccount(b)
	common.FatalIfError(t, err, "inserting account")

	as, err = db.SelectAccounts()
	common.FatalIfError(t, err, "selecting accounts after inserting two")

	if !assert.Len(t, *as, 2) {
		t.FailNow()
	}
	retrievedB := (*as)[1]
	equal, err = insertedB.Equal(retrievedB)
	common.FatalIfError(t, err, "equaling inserted and retrieved")
	if !assert.True(t, equal) {
		t.FailNow()
	}

	equal, err = insertedA.Equal(*insertedB)
	common.FatalIfError(t, err, "equaling insertedA and insertedB")
	if !assert.False(t, equal) {
		t.FailNow()
	}

	equal, err = retrievedA.Equal(retrievedB)
	common.FatalIfError(t, err, "equaling retrievedA and retrievedB")
	if !assert.False(t, equal) {
		t.FailNow()
	}
}

func createStorage(t *testing.T) storage.Storage {
	cs, err := postgres2.NewConnectionString(
		viper.GetString(keyDBHost),
		viper.GetString(keyDBUser),
		viper.GetString(keyDBName),
		viper.GetString(keyDBSSLMode),
	)
	common.FatalIfError(t, err, "creating connection string")
	db, err := postgres2.New(cs)
	common.FatalIfError(t, err, "creating storage")
	return db
}
