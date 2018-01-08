package main

import (
	"log"

	"github.com/glynternet/accounting-rest/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("error executing root command: %v", err)
	}
}
