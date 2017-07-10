package main

import (
	"log"
	"net/http"
	"fmt"
	"errors"
	"os"
	"io"
)

var connectionString string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No database connection file location given. Please provide the location of the connection string file as the first argument to the application.")
		return
	}
	var err error
	connectionString, err = loadDBConnectionString(os.Args[1])
	if err != nil {
		fmt.Printf("Unable to load connection string from file at %s\n", os.Args[1])
		return
	}
	router := NewRouter()
	port := 8080
	fmt.Printf("Starting GOHMoneyREST on port %d\n", port)
	portString := fmt.Sprintf(`:%d`, port)
	log.Fatal(http.ListenAndServe(portString, router))
}

// loadDBConnectionString loads the connection string to be used when attempting to connect to the database throught the application.
func loadDBConnectionString(connectionStringLocation string) (string, error) {
	if len(connectionStringLocation) < 1 {
		return ``, errors.New("No connection string file location given.")
	}
	file, err := os.Open(connectionStringLocation)
	if err != nil {
		return ``, err
	}
	stat, err := file.Stat()
	if err != nil {
		return ``, err
	}
	maxConnectionString := int64(200)
	fileSize := stat.Size()
	if fileSize > maxConnectionString {
		message := fmt.Sprintf("Connection string file (%s) is too large. Max: %d, Length: %d", connectionStringLocation, maxConnectionString, fileSize)
		return ``, errors.New(message)
	}
	connectionString := make([]byte, maxConnectionString)
	bytesCount, err := file.Read(connectionString)
	if err != nil && err != io.EOF {
		return ``, err
	}
	return string(connectionString[0:bytesCount]), err
}

