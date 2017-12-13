package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

func Accounts(w http.ResponseWriter, r *http.Request) (int, error) {
	store, err := NewStorage()
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	as, err := store.SelectAccounts()
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
	if err := json.NewEncoder(w).Encode(as); err != nil {
		log.Printf("error encoding json: %v", err)
	}
	return http.StatusOK, nil
}
