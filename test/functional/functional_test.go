package functional

import (
	"strings"
	"testing"

	"github.com/glynternet/go-accounting-storage/postgres2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	appName = "accounting-rest"

	// viper keys
	keyDBHost    = "db-host"
	keyDBUser    = "db-user"
	keyDBName    = "db-name"
	keyDBSSLMode = "db-sslmode"
)

func TestInit(t *testing.T) {
	err := postgres2.CreateStorage(
		viper.GetString(keyDBHost),
		viper.GetString(keyDBUser),
		viper.GetString(keyDBName),
		viper.GetString(keyDBSSLMode),
	)
	assert.NoError(t, err)
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
