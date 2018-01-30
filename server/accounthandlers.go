package server

import (
	"encoding/json"
	"net/http"

	"github.com/glynternet/go-accounting-storage"
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
	w.WriteHeader(http.StatusOK)
	w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
	return http.StatusOK, errors.Wrap(
		json.NewEncoder(w).Encode(as),
		"error encoding accounts json",
	)
}
