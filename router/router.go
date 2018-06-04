package router

import (
	"log"
	"net/http"

	"github.com/glynternet/mon/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// New creates a new mux.Router and initialises it with routes for the store
func New(store storage.Storage) (*mux.Router, error) {
	if store == nil {
		return nil, errors.New("nil store")
	}
	return (&router{storage: store}).New()
}

// New creates a new Router and initialises it will all of the global routes
func (s router) New() (*mux.Router, error) {
	// TODO: can this receiver just be a value not a pointer
	//if s == nil {
	//	return nil, errors.New("nil router provided")
	//}
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range s.routes() {
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

type router struct {
	storage storage.Storage
}

func (s router) routes() []route {
	return []route{
		{
			name:       "Accounts",
			pattern:    patternAccounts,
			appHandler: s.handlerSelectAccounts,
			method:     http.MethodGet,
		},
		{
			name:       "Account",
			pattern:    patternAccount,
			appHandler: s.muxAccountIDHandlerFunc,
			method:     http.MethodGet,
		},
		{
			name:       "AccountInsert",
			pattern:    EndpointAccountInsert,
			appHandler: s.muxAccountInsertHandlerFunc,
			method:     http.MethodPost,
		},
		{
			name:       "AccountUpdate",
			pattern:    patternAccountUpdate,
			appHandler: s.muxAccountUpdateHandlerFunc,
			method:     http.MethodPost,
		},
		{
			name:       "Balances",
			pattern:    patternAccountBalances,
			appHandler: s.muxAccountBalancesHandlerFunc,
			method:     http.MethodGet,
		},
		{
			name:       "BalanceInsert",
			pattern:    patternAccountBalanceInsert,
			appHandler: s.muxAccountBalanceInsertHandlerFunc,
			method:     http.MethodPost,
		},
	}
}
