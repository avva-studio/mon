package router

import (
	"log"
	"net/http"
	"time"
)

// logger returns a simple logger that wraps a http.handler and logs any incoming request
func logger(inner appJSONHandler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
