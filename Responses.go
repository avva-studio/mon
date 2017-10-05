package main

import (
	"net/http"
	"log"
)

// ServiceUnavailableResponse writes to the response header and body that the GOHMoneyREST service is unavailable.
func ServiceUnavailableResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusServiceUnavailable)
	_, err := w.Write([]byte(`Service currently unavailable.`))
	if err != nil {
		log.Printf("Error writing response: %s", err)
	}

}
