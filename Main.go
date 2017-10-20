package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/glynternet/GOHMoneyDB"
)

var connectionString string

func main() {
	router := NewRouter()
	port := 8080
	log.Printf("Starting GOHMoneyREST on port %d\n", port)
	portString := fmt.Sprintf(`:%d`, port)
	log.Fatal(http.ListenAndServe(portString, router))
}

func init() {
	if len(os.Args) < 2 {
		log.Println("No database connection file location given. Please provide the location of the connection string file as the first argument to the application.")
		return
	}
	var err error
	connectionString, err = GOHMoneyDB.LoadDBConnectionString(os.Args[1])
	if err != nil {
		log.Printf("Unable to load connection string from file at %s\n", os.Args[1])
		return
	}
	logDBState()
}

func logDBState() {
	db, err := GOHMoneyDB.OpenDBConnection(connectionString)
	if err != nil {
		log.Printf("Error opening database connection: %s", err)
		return
	}
	defer db.Close()
	msg := "Database is "
	if !GOHMoneyDB.DbIsAvailable(db) {
		msg += "not "
	}
	log.Print(msg + "available.")
	return
}