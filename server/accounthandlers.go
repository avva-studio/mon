package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/glynternet/go-accounting-storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) accounts(w http.ResponseWriter, _ *http.Request) (int, error) {
	if w == nil {
		return http.StatusInternalServerError, errors.New("nil ResponseWriter")
	}
	store, err := s.NewStorage()
	if err != nil {
		return http.StatusServiceUnavailable, errors.Wrap(err, "creating new storage")
	}
	var as *storage.Accounts
	as, err = store.SelectAccounts()
	if err != nil {
		return http.StatusServiceUnavailable, errors.Wrap(err, "selecting accounts from client")
	}
	w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
	return http.StatusOK, errors.Wrap(
		json.NewEncoder(w).Encode(as),
		"error encoding accounts json",
	)
}

func (s *server) muxAccountIDHandlerFunc(w http.ResponseWriter, r *http.Request) (int, error) {
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
	return s.account(uint(id))(w, r)
}

func (s *server) account(id uint) appHandler {
	return func(w http.ResponseWriter, r *http.Request) (int, error) {
		if w == nil {
			return http.StatusInternalServerError, errors.New("nil ResponseWriter")
		}
		store, err := s.NewStorage()
		if err != nil {
			return http.StatusServiceUnavailable, errors.Wrap(err, "creating new storage")
		}
		var a *storage.Account
		a, err = store.SelectAccount(id)
		if err != nil {
			return http.StatusNotFound, errors.Wrap(err, "selecting account from storage")
		}
		w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
		return http.StatusOK, errors.Wrap(
			json.NewEncoder(w).Encode(a),
			"error encoding accounts json",
		)
	}
}
