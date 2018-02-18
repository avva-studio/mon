package server

import (
	"net/http"
	"strconv"

	"github.com/glynternet/go-accounting-storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) handlerSelectAccounts(_ *http.Request) (int, interface{}, error) {
	store, err := s.NewStorage()
	if err != nil {
		return http.StatusServiceUnavailable, nil, errors.Wrap(err, "creating new storage")
	}
	var as *storage.Accounts
	as, err = store.SelectAccounts()
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
	return s.handlerSelectAccount(uint(id))(r)
}

func (s *server) handlerSelectAccount(id uint) appJSONHandler {
	return func(_ *http.Request) (int, interface{}, error) {
		store, err := s.NewStorage()
		if err != nil {
			return http.StatusServiceUnavailable, nil, errors.Wrap(err, "creating new storage")
		}
		var a *storage.Account
		a, err = store.SelectAccount(id)
		if err != nil {
			return http.StatusNotFound, nil, errors.Wrap(err, "selecting handlerSelectAccount from storage")
		}
		return http.StatusOK, a, nil
	}
}

///// THIS NEEDS SPLITTING UP AND ISN'T GOING SO WELL AT THE MOMENT
///// I THINK IT WOULD BE BEST TO SPLIT UP INTO SOMETHING THAT TAKES IN A REQUEST
///// AND PARSES IT INTO AN ACCOUNT, THEN WE CAN JUST HAVE A HANDLER THAT TAKES
///// AN ACCOUNT AND WRITES TO A RESPONSE WRITER???

//func (s *server) handlerInsertAccout2(w http.ResponseWriter, r *http.Request) {
//	a, err := parseAccount(r)
//	sa, err := insertAccount(a)
//
//}

//func insertAccountStorageRequestHandlerfunc(store storage.Storage, r *http.Request) (interface{}, int, error) {
//	inner, err := readAccount(r.Body)
//	if err != nil {
//		return nil, http.StatusBadRequest, errors.Wrap(err, "reading account from body")
//	}
//	a, err := store.InsertAccount(*inner)
//	if err != nil {
//		return nil, http.StatusBadRequest, errors.Wrap(err, "inserting account to storage")
//	}
//	return a, http.StatusOK, nil
//}
//
//func (s *server) newHandlerInsertAccountAppHandler() appJSONHandler {
//	store, err := s.NewStorage()
//	if err != nil {
//		return nil
//	}
//	return func(w http.ResponseWriter, r *http.Request) (int, error) {
//		insertAccountStorageRequestHandlerfunc(store, r)
//	}
//}
//
//func (s *server) handlerInsertAccount(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
//
//	if w == nil {
//		return http.StatusInternalServerError, nil, errors.New("nil ResponseWriter")
//	}
//
//	store, err := s.NewStorage()
//	if err != nil {
//		return http.StatusServiceUnavailable, nil, errors.Wrap(err, "creating new storage")
//	}
//
//	encode, code, err := func(store storage.Storage, r *http.Request) (interface{}, int, error) {
//		inner, err := readAccount(r.Body)
//		if err != nil {
//			return nil, http.StatusBadRequest, errors.Wrap(err, "reading account from body")
//		}
//		a, err := store.InsertAccount(*inner)
//		if err != nil {
//			return nil, http.StatusBadRequest, errors.Wrap(err, "inserting account to storage")
//		}
//		return a, http.StatusOK, nil
//	}(store, r)
//
//	if err != nil {
//		return code, errors.Wrap(err, "handling insert account request")
//	}
//
//	return http.StatusOK, encode, nil
//}
//
//func readAccount(r io.Reader) (*account2.Account, error) {
//	bod, err := ioutil.ReadAll(r)
//	if err != nil {
//		return nil, errors.Wrap(err, "reading all body")
//	}
//	inner := new(account2.Account)
//	return inner, errors.Wrap(
//		json.Unmarshal(bod, inner),
//		"unmarshalling body to account",
//	)
//}
