package server

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// newRouter creates a new Router and initialises it will all of the global routes
func (s *server) newRouter() (*mux.Router, error) {
	if s == nil {
		return nil, errors.New("nil server provided")
	}
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range s.routes() {
		handler := logger(route.appHandler, route.name)
		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			Handler(handler)
	}
	return router, nil
}
