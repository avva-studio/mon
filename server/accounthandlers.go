package server

import (
	"net/http"
	"strconv"

	"github.com/glynternet/go-accounting/account"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// TODO: redesign these so that they don't need to take a request? There could be multiple handler types either take a request or don't take a request
func (s *server) handlerSelectAccounts(_ *http.Request) (int, interface{}, error) {
	as, err := s.storage.SelectAccounts()
	if err != nil {
		return http.StatusServiceUnavailable, nil, errors.Wrap(err, "selecting handlerSelectAccounts from client")
	}
	return http.StatusOK, as, nil
}

func (s *server) muxAccountIDHandlerFunc(r *http.Request) (int, interface{}, error) {
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
	return s.handlerSelectAccount(uint(id))
}

func (s *server) handlerSelectAccount(id uint) (int, interface{}, error) {
	a, err := s.storage.SelectAccount(id)
	if err != nil {
		return http.StatusNotFound, nil, errors.Wrap(err, "selecting handlerSelectAccount from storage")
	}
	return http.StatusOK, a, nil
}

func (s *server) handlerInsertAccount(a account.Account) (int, interface{}, error) {
	return nil, nil, nil
}
