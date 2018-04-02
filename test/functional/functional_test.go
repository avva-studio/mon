package functional

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/glynternet/accounting-rest/client"
	"github.com/glynternet/go-accounting-storage/test"
	"github.com/spf13/viper"
)

const (
	// viper keys
	keyServerHost = "server-host"
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
	setup()
}

func setup() {
	const retries = 5
	errs := make([]error, retries)
	store := client.Client(viper.GetString(keyServerHost))
	var i int
	for i = 0; i < retries; i++ {
		_, err := store.SelectAccounts()
		if err == nil {
			break
		}
		errs[i] = err
		time.Sleep(time.Second)
	}
	if errs[retries-1] != nil {
		for i, err := range errs {
			fmt.Printf("[retry: %02d] %v\n", i, err)
		}
	}

}

func TestSuite(t *testing.T) {
	store := client.Client(viper.GetString(keyServerHost))
	if !store.Available() {
		t.Fatal("store is unavailable")
	}
	test.Test(t, store)
}
