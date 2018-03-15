package functional

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/glynternet/go-accounting-storage/postgres2"
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

func TestInit(t *testing.T) {
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
		log.Printf("Attempt: %d, err: %v\n", i, err)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	assert.NoError(t, err)
	t.Logf("Attempts: %d", i)
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
}

//func newStorage(host, user, dbname, sslmode string) (storage.Storage, error) {
//	cs, err := postgres2.NewConnectionString(host, user, dbname, sslmode)
//	if err != nil {
//		return nil, fmt.Errorf("unable to create connection string: %v", err)
//	}
//	return postgres2.New(cs)
//}
