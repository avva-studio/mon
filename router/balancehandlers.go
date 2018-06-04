package router

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *router) balances(accountID uint) (int, interface{}, error) {
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

func (s *router) muxAccountBalancesHandlerFunc(r *http.Request) (int, interface{}, error) {
	id, err := extractID(mux.Vars(r))
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "extracting account ID")
	}
	return s.balances(id)
}

func (s *router) insertBalance(accountID uint, b balance.Balance) (int, interface{}, error) {
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

func (s *router) muxAccountBalanceInsertHandlerFunc(r *http.Request) (int, interface{}, error) {
	id, err := extractID(mux.Vars(r))
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "extracting account ID")
	}

	bod, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "reading request body")
	}

	defer func() {
		// TODO: this handler only needs to take a []byte which would mean we can handle closing the body elsewhere
		cErr := r.Body.Close()
		if cErr != nil {
			log.Print(errors.Wrap(err, "closing request body"))
		}
	}()

	var b balance.Balance
	err = json.Unmarshal(bod, &b)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "unmarshalling request body")
	}
	return s.insertBalance(id, b)
}
