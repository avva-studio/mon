package internal

import (
	"github.com/glynternet/go-accounting-storage"
	"github.com/gorilla/mux"
)

var NewStorage StorageFunc

type StorageFunc func() (storage.Storage, error)

// NewRouter creates a new Router and initialises it will all of the global routes
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := logger(route.appHandler)
		router.
			Methods(route.appHandler.method).
			Path(route.pattern).
			Name(route.appHandler.name).
			Handler(handler)
	}
	return router
}
