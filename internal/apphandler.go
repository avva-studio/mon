package internal

import (
	"log"
	"net/http"
)

type appHandler struct {
	name    string
	method  string
	handler func(http.ResponseWriter, *http.Request) (int, error)
}

// ServeHTTP makes our appHandler function satisfy the http.HandlerFunc interface
// If we are returning an error from our appHandler, we should not have already
// written to our ResponseWriter
func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := ah.handler(w, r); err != nil {
		log.Printf(
			"error serving on appHandler %v. Error: %v - Status: %d (%s) - Request: %+v",
			ah, err, status, http.StatusText(status), r,
		)
		switch status {
		case http.StatusServiceUnavailable:
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			// We can have cases as granular as we like, if we wanted to
			// return custom errors for specific status codes.
			// TODO: if http.StatusInternalServerError is received, we should return bad request and log the error maybe?
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			// Catch any other errors we haven't explicitly handled
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
