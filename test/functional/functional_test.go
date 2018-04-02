package functional

import (
	"strings"
	"testing"

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

func setup() {}

func TestSuite(t *testing.T) {
	store := client.Client(viper.GetString(keyServerHost))
	if !store.Available() {
		t.Fatal("store is unavailable")
	}
	test.Test(t, store)
}
