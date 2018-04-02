package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/glynternet/accounting-rest/server"
	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storage/postgres2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName = "accounting-rest-serve"

	// viper keys
	keyPort      = "port"
	keyDBHost    = "db-host"
	keyDBUser    = "db-user"
	keyDBName    = "db-name"
	keyDBSSLMode = "db-sslmode"
)

func main() {
	log.Fatal(cmdDBServe.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	cmdDBServe.Flags().String(keyPort, "8080", "server listening port")
	cmdDBServe.Flags().String(keyDBHost, "", "host address of the DB backend")
	cmdDBServe.Flags().String(keyDBName, "", "name of the DB set to use")
	cmdDBServe.Flags().String(keyDBUser, "", "DB user to authenticate with")
	cmdDBServe.Flags().String(keyDBSSLMode, "", "DB SSL mode to use")
	err := viper.BindPFlags(cmdDBServe.Flags())
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

var cmdDBServe = &cobra.Command{
	Use: appName,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print(viper.GetString(keyDBHost),
			viper.GetString(keyDBUser),
			viper.GetString(keyDBName),
			viper.GetString(keyDBSSLMode))
		store, err := newStorage(
			viper.GetString(keyDBHost),
			viper.GetString(keyDBUser),
			viper.GetString(keyDBName),
			viper.GetString(keyDBSSLMode),
		)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error creating storage"))
		}
		s, err := server.New(store)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error creating new server"))
		}
		log.Fatal(s.ListenAndServe(":" + viper.GetString(keyPort)))
	},
}
