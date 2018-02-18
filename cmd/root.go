package cmd

import (
	"fmt"
	"log"

	"github.com/glynternet/accounting-rest/server"
	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storage/postgres2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	host    string
	user    string
	dbname  string
	sslmode string
	port    string
)

const (
	keyHost    = "host"
	keyUser    = "user"
	keyDBName  = "dbname"
	keySSLMode = "sslmode"
	appName    = "accounting-rest"
)

func Execute() error {
	return RootCmd.Execute()
}

var RootCmd = &cobra.Command{
	Use: appName,
	Run: func(cmd *cobra.Command, args []string) {
		host = viper.GetString(keyHost)
		log.Printf("%s %s", keyHost, host)
		user = viper.GetString(keyUser)
		log.Printf("%s %s", keyUser, user)
		dbname = viper.GetString(keyDBName)
		log.Printf("%s %s", keyDBName, dbname)
		sslmode = viper.GetString(keySSLMode)
		log.Printf("%s %s", keySSLMode, sslmode)
		store, err := newStorage(host, user, dbname, sslmode)
		if err != nil {
			log.Fatalf("error creating storage: %v", err)
		}
		s, err := server.New(store)
		if err != nil {
			log.Fatalf("error creating new server")
		}
		log.Fatal(s.ListenAndServe(":" + port))
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&port, "port", "8080", "server listening port")
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/." + appName)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
}

func newStorage(host, user, dbname, sslmode string) (storage.Storage, error) {
	cs, err := postgres2.NewConnectionString(host, user, dbname, sslmode)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection string: %v", err)
	}
	return postgres2.New(cs)
}
