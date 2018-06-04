package router

import (
	"log"
	"net/http"

	"github.com/glynternet/mon/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// New creates a new mux.Router and initialises it with generateRoutes for the store
func New(store storage.Storage) (*mux.Router, error) {
	if store == nil {
		return nil, errors.New("nil store")
	}
	rs := generateRoutes(environment{storage: store})
	return newRouter(rs)
}

// New creates a new Router and initialises it will all of the global generateRoutes
func newRouter(rs []route) (*mux.Router, error) {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range rs {
		handler := logger(route.appHandler, route.name)
		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			Handler(handler)
		log.Printf("Route registered: %+v", route)
	}
	return router, nil
}

type environment struct {
	storage storage.Storage
}

func generateRoutes(e environment) []route {
	return []route{
		{
			name:       "Accounts",
			pattern:    patternAccounts,
			appHandler: e.handlerSelectAccounts,
			method:     http.MethodGet,
		},
		{
			name:       "Account",
			pattern:    patternAccount,
			appHandler: e.muxAccountIDHandlerFunc,
			method:     http.MethodGet,
		},
		{
			name:       "AccountInsert",
			pattern:    EndpointAccountInsert,
			appHandler: e.muxAccountInsertHandlerFunc,
			method:     http.MethodPost,
		},
		{
			name:       "AccountUpdate",
			pattern:    patternAccountUpdate,
			appHandler: e.muxAccountUpdateHandlerFunc,
			method:     http.MethodPost,
		},
		{
			name:       "Balances",
			pattern:    patternAccountBalances,
			appHandler: e.muxAccountBalancesHandlerFunc,
			method:     http.MethodGet,
		},
		{
			name:       "BalanceInsert",
			pattern:    patternAccountBalanceInsert,
			appHandler: e.muxAccountBalanceInsertHandlerFunc,
			method:     http.MethodPost,
		},
	}
}
