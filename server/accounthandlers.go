package server

import (
	"net/http"
	"strconv"

	"io/ioutil"

	"log"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// TODO: redesign these so that they don't need to take a request? There could be multiple handler types either take a request or don't take a request
func (s *server) handlerSelectAccounts(_ *http.Request) (int, interface{}, error) {
	as, err := s.storage.SelectAccounts()
	if err != nil {
		return http.StatusServiceUnavailable, nil, errors.Wrap(err, "selecting Accounts from client")
	}
	return http.StatusOK, as, nil
}

func (s *server) muxAccountIDHandlerFunc(r *http.Request) (int, interface{}, error) {
	id, err := extractID(mux.Vars(r))
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "extracting account ID")
	}
	return s.handlerSelectAccount(id)
}

func (s *server) handlerSelectAccount(id uint) (int, interface{}, error) {
	a, err := s.storage.SelectAccount(id)
	if err != nil {
		return http.StatusNotFound, nil, errors.Wrapf(err, "selecting Account with id:%d from storage", id)
	}
	return http.StatusOK, a, nil
}

func (s *server) muxAccountInsertHandlerFunc(r *http.Request) (int, interface{}, error) {
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

	a, err := account.UnmarshalJSON(bod)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "unmarshalling request body")
	}
	return s.handlerInsertAccount(*a)
}

func (s *server) handlerInsertAccount(a account.Account) (int, interface{}, error) {
	inserted, err := s.storage.InsertAccount(a)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, "inserting Account into storage")
	}
	return http.StatusOK, inserted, nil
}

func (s *server) muxAccountUpdateHandlerFunc(r *http.Request) (int, interface{}, error) {
	id, err := extractID(mux.Vars(r))
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "extracting account ID")
	}

	o, err := s.storage.SelectAccount(id)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "selecting account with id:%d", id)
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

	updates, err := account.UnmarshalJSON(bod)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "unmarshalling request body")
	}

	return s.handlerUpdateAccount(*o, *updates)
}

func (s *server) handlerUpdateAccount(original storage.Account, updates account.Account) (int, interface{}, error) {
	updated, err := s.storage.UpdateAccount(&original, &updates)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	return http.StatusOK, updated, nil
}

func (s *server) handlerDeleteAccount(id uint) (int, interface{}, error) {
	err := s.storage.DeleteAccount(id)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "deleting Account with id:%d from storage", id)
	}
	return http.StatusOK, nil, nil
}

func extractID(vars map[string]string) (uint, error) {
	if vars == nil {
		return 0, errors.New("nil vars map")
	}
	idString, ok := vars["id"]
	if !ok {
		return 0, errors.New("no account id context variable")
	}
	id, err := strconv.ParseUint(idString, 10, 64)
	return uint(id), errors.Wrapf(err, "parsing %s to uint", idString)
}
