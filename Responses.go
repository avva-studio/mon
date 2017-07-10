package main

import (
	"net/http"
)

// ServiceUnavailableResponse writes to the response header and body that the GOHMoneyREST service is unavailable.
func ServiceUnavailableResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(`Service currently unavailable.`))
}
