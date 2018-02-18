package server

import (
	"net/http"
	"strconv"

	"github.com/glynternet/go-accounting-storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) balances(accountId uint) appJSONHandler {
	return func(r *http.Request) (int, interface{}, error) {
		store, err := s.NewStorage()
		if err != nil {
			return http.StatusServiceUnavailable, nil, errors.Wrap(err, "creating new storage")
		}
		a, err := store.SelectAccount(accountId)
		if err != nil {
			return http.StatusBadRequest, nil, errors.Wrapf(err, "selecting account with id %d", accountId)
		}
		var bs *storage.Balances
		bs, err = store.SelectAccountBalances(*a)
		if err != nil {
			return http.StatusBadRequest, nil, errors.Wrapf(err, "selecting balances for account %+v", *a)
		}
		return http.StatusOK, bs, nil
	}
}

func (s *server) muxAccountBalancesHandlerFunc(r *http.Request) (int, interface{}, error) {
	vars := mux.Vars(r)
	if vars == nil {
		return http.StatusBadRequest, nil, errors.New("no context variables")
	}

	key := "id"
	idString, ok := vars[key]
	if !ok {
		return http.StatusBadRequest, nil, errors.New("no account_id context variable")
	}
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "parsing %s to uint", key)
	}
	return s.balances(uint(id))(r)
}
