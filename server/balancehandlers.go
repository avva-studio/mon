package server

import (
	"net/http"

	"github.com/pkg/errors"
)

func (s *server) balances(accountId uint) appHandler {
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
		_, err = store.SelectAccountBalances(*a)
		if err != nil {
			return http.StatusBadRequest, errors.Wrapf(err, "selecting balances for account %+v", *a)
		}
		return 0, errors.New("not implemented")
	}
}
