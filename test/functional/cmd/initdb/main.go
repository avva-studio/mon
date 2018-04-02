package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/glynternet/go-accounting-storage/postgres2"
	"github.com/spf13/viper"
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

func main() {
	const retries = 5
	errs := make([]error, retries)
	var i int
	for i = 0; i < retries; i++ {
		err := postgres2.CreateStorage(
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
			fmt.Printf("[retry: %02d] %v\n", i, err)
		}
		os.Exit(1)
	}
}
