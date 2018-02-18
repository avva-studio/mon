package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/glynternet/go-accounting-storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) balances(accountId uint) appJSONHandler {
	return func(w http.ResponseWriter, r *http.Request) (int, error) {
		if w == nil {
			return http.StatusInternalServerError, errors.New("nil ResponseWriter")
		}
		store, err := s.NewStorage()
		if err != nil {
			return http.StatusServiceUnavailable, errors.Wrap(err, "creating new storage")
		}
		a, err := store.SelectAccount(accountId)
		if err != nil {
			return http.StatusBadRequest, errors.Wrapf(err, "selecting account with id %d", accountId)
		}
		var bs *storage.Balances
		bs, err = store.SelectAccountBalances(*a)
		if err != nil {
			return http.StatusBadRequest, errors.Wrapf(err, "selecting balances for account %+v", *a)
		}
		w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
		return http.StatusOK, errors.Wrap(
			json.NewEncoder(w).Encode(bs),
			"error encoding balances json",
		)
	}
}

func (s *server) muxAccountBalancesHandlerFunc(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	if vars == nil {
		return http.StatusBadRequest, errors.New("no context variables")
	}

	key := "id"
	idString, ok := vars[key]
	if !ok {
		return http.StatusBadRequest, errors.New("no account_id context variable")
	}
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		return http.StatusBadRequest, errors.Wrapf(err, "parsing %s to uint", key)
	}
	return s.balances(uint(id))(w, r)
}
