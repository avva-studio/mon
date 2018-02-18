package server

import (
	"net/http"
	"strconv"

	"github.com/glynternet/go-accounting-storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) balances(accountID uint) appJSONHandler {
	return func(_ *http.Request) (int, interface{}, error) {
		a, err := s.storage.SelectAccount(accountID)
		if err != nil {
			return http.StatusBadRequest, nil, errors.Wrapf(err, "selecting account with id %d", accountID)
		}
		var bs *storage.Balances
		bs, err = s.storage.SelectAccountBalances(*a)
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
	return s.balances(uint(id))(nil) // Request is not needed for balances handler
}
