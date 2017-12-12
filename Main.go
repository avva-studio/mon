package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/glynternet/go-accounting-storage/postgres2"
	"github.com/spf13/viper"
)

var connectionString string

var (
	host    string
	user    string
	dbname  string
	sslmode string
)

const (
	keyHost    = "host"
	keyUser    = "user"
	keyDBName  = "dbname"
	keySSLMode = "sslmode"
)

func main() {

	router := NewRouter()
	port := 8080
	log.Printf("Starting GOHMoneyREST on port %d\n", port)
	portString := fmt.Sprintf(`:%d`, port)
	log.Fatal(http.ListenAndServe(portString, router))
}

func init() {
	viper.SetConfigName("config")                 // name of config file (without extension)
	viper.AddConfigPath("$HOME/.accounting-rest") // call multiple times to add many search paths
	viper.AddConfigPath(".")                      // optionally look for config in the working directory
	err := viper.ReadInConfig()                   // Find and read the config file
	if err != nil {                               // Handle errors reading the config file
		log.Fatalf("error reading config: %v", err)
	}
	host = viper.GetString(keyHost)
	log.Printf("%s %s", keyHost, host)
	user = viper.GetString(keyUser)
	log.Printf("%s %s", keyUser, user)
	dbname = viper.GetString(keyDBName)
	log.Printf("%s %s", keyDBName, dbname)
	sslmode = viper.GetString(keySSLMode)
	log.Printf("%s %s", keySSLMode, sslmode)

	os.Exit(0)

	logDBState()

}

func logDBState() {
	cs, err := postgres2.NewConnectionString(host, user, dbname, sslmode)
	if err != nil {
		log.Fatalf("unable to create connection string: %v", err)
	}
	pg, err := postgres2.New(cs)
	if err != nil {
		log.Fatalf("unable to create postgres store: %v", err)
	}
	defer pg.Close()
	msg := "Storage is "
	if !pg.Available() {
		msg += "not "
	}
	log.Print(msg + "available.")
	return
}
