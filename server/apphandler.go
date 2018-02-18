package server

import (
	"encoding/json"
	"log"
	"net/http"

	"io"

	"github.com/pkg/errors"
)

type appJSONHandler func(http.ResponseWriter, *http.Request) (int, interface{}, error)

// ServeHTTP makes our appJSONHandler function satisfy the http.HandlerFunc interface
// If we are returning an error from our appJSONHandler, we should not have already
// written to our ResponseWriter
func (ah appJSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, bod, err := ah(w, r)

	// handle errors
	if err != nil {
		log.Printf(
			"error serving on appJSONHandler %v. Error: %v - Status: %d (%s) - Request: %+v",
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
		return
	}

	// here, I don't want to write to the writer immediately using a json
	// encoder, in case there is an error in json encoding
	bs, err := json.Marshal(bod)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, wErr := io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		if wErr != nil {
			log.Print(errors.Wrap(wErr, "writing status text to ResponseWriter"))
		}
		log.Print(errors.Wrap(err, "marshalling json reponse"))
		return
	}

	w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
	_, wErr := w.Write(bs)
	if wErr != nil {
		log.Print(errors.Wrap(wErr, "writing body to ResponseWriter"))
	}
}
