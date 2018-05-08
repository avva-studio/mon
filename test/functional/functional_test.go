// +build functional

package functional

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/glynternet/accounting-rest/client"
	"github.com/glynternet/go-accounting-storage/postgres"
	"github.com/glynternet/go-accounting-storage/test"
	"github.com/spf13/viper"
)

const (
	// viper keys
	keyServerHost = "server-host"
	keyDBHost     = "db-host"
	keyDBUser     = "db-user"
	keyDBName     = "db-name"
	keyDBSSLMode  = "db-sslmode"
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
	setup()
}

func setup() {
	const retries = 5
	errs := make([]error, retries)
	var i int
	for i = 0; i < retries; i++ {
		err := postgres.CreateStorage(
			viper.GetString(keyDBHost),
			viper.GetString(keyDBUser),
			viper.GetString(keyDBName),
			viper.GetString(keyDBSSLMode),
		)
		if err == nil {
			break
		}
		errs[i] = err
		time.Sleep(time.Second)
	}
	if errs[retries-1] != nil {
		for i, err := range errs {
			log.Printf("[retry: %02d] %v\n", i, err)
		}
		os.Exit(1)
	}
	log.Print("Setup complete")
}

func TestSuite(t *testing.T) {
	store := client.Client(viper.GetString(keyServerHost))
	if !store.Available() {
		t.Fatal("store is unavailable")
	}
	test.Test(t, store)
}
