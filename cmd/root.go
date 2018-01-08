package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/glynternet/accounting-rest/internal"
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
		storageFn, err := newStorageFunc(host, user, dbname, sslmode)
		if err != nil {
			log.Fatalf("error creating storage func: %v", err)
		}
		internal.NewStorage = storageFn
		logDBState()
		router := internal.NewRouter()
		log.Printf("Starting %s on port %s\n", appName, port)
		log.Fatal(http.ListenAndServe(":"+port, router))
	},
}

func logDBState() {
	store, err := internal.NewStorage()
	if err != nil {
		log.Printf("error creating new storage: %v", err)
	}
	defer store.Close()
	msg := "Storage is "
	if !store.Available() {
		msg += "not "
	}
	log.Print(msg + "available.")
	return
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

func newStorageFunc(host, user, dbname, sslmode string) (internal.StorageFunc, error) {
	cs, err := postgres2.NewConnectionString(host, user, dbname, sslmode)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection string: %v", err)
	}
	return func() (storage.Storage, error) {
		return postgres2.New(cs)
	}, nil
}
