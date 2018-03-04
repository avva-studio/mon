package server

import (
	"net/http"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting/balance"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) balances(accountID uint) (int, interface{}, error) {
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

func (s *server) muxAccountBalancesHandlerFunc(r *http.Request) (int, interface{}, error) {
	id, err := extractID(mux.Vars(r))
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "extracting account ID")
	}
	return s.balances(uint(id))
}

func (s *server) insertBalance(accountID uint, b balance.Balance) (int, interface{}, error) {
	a, err := s.storage.SelectAccount(accountID)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, "selecting account")
	}
	inserted, err := s.storage.InsertBalance(*a, b)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, "inserting balance")
	}
	return http.StatusOK, inserted, nil
}
