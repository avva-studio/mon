package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storage/postgres2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName = "accounting-rest"

	// viper keys
	keyPort      = "port"
	keyDBHost    = "db-host"
	keyDBUser    = "db-user"
	keyDBName    = "db-name"
	keyDBSSLMode = "db-sslmode"
)

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use: appName,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String(keyPort, "8080", "server listening port")
	rootCmd.PersistentFlags().String(keyDBHost, "", "host address of the DB backend")
	rootCmd.PersistentFlags().String(keyDBName, "", "name of the DB set to use")
	rootCmd.PersistentFlags().String(keyDBUser, "", "DB user to authenticate with")
	rootCmd.PersistentFlags().String(keyDBSSLMode, "", "DB SSL mode to use")
	err := viper.BindPFlags(rootCmd.Flags())
	if err != nil {
		log.Printf("unable to BindPFlags: %v", err)
	}
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
}

func newStorage(host, user, dbname, sslmode string) (storage.Storage, error) {
	cs, err := postgres2.NewConnectionString(host, user, dbname, sslmode)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection string: %v", err)
	}
	return postgres2.New(cs)
}
