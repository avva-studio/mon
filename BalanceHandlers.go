package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/glynternet/GOHMoney/balance"
	"github.com/glynternet/GOHMoneyDB"
	"github.com/gorilla/mux"
	"log"
)


// BalanceCreate handler accepts json representing a potential new GOHMoney.Balance. The Balance is decoded and attempted to be added to the backend.
// If successful, the response contains json representing the newly created GOHMoneyDB.Balance object,
// else, an error describing why the creation was unsuccessful.
func BalanceCreate(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
	var newBalance accountBalanceJSONHelper
	err := decoder.Decode(&newBalance)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, werr := w.Write([]byte("Error decoding new account: " + err.Error()))
		if werr != nil {
			log.Printf("Error writing to bytes: %s", werr)
		}
		return
	}
	defer r.Body.Close()
	if newBalance.AccountID < 1 {
		w.WriteHeader(http.StatusBadRequest)
		_, werr := w.Write([]byte(`Account ID must be a positive integer`))
		if werr != nil {
			log.Printf("Error writing to bytes: %s", werr)
		}
		return
	}
	db, err := GOHMoneyDB.OpenDBConnection(connectionString)
	if err != nil {
		ServiceUnavailableResponse(w)
		return
	}
	defer db.Close()
	if !GOHMoneyDB.DbIsAvailable(db) {
		ServiceUnavailableResponse(w)
		return
	}
	account, err := GOHMoneyDB.SelectAccountWithID(db, newBalance.AccountID)
	if err != nil {
		if _, ok := err.(GOHMoneyDB.NoAccountWithIDError); ok {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}
	if err := newBalance.Balance.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	createdBalance, err := GOHMoneyDB.Account(account).InsertBalance(db, newBalance.Balance)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	balanceData, err := json.Marshal(createdBalance)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error creating json from created balance data: %s", err.Error())))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(balanceData)
}

// BalanceUpdate handler accepts json representing a potential update to a GOHMoney.Balance object along with the id of the account owner. The Balance is decoded and attempted to be updated in the backend.
// If successful, the response contains json representing the newly updated GOHMoneyDB.Balance object and returns a 204 status.
// else, an error describing why the update was unsuccessful.
func BalanceUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseIDString(vars[`id`])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error parsing id: %s", err.Error())))
		return
	}
	db, err := GOHMoneyDB.OpenDBConnection(connectionString)
	if err != nil {
		ServiceUnavailableResponse(w)
		return
	}
	defer db.Close()
	if !GOHMoneyDB.DbIsAvailable(db) {
		ServiceUnavailableResponse(w)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var newBalance accountBalanceJSONHelper
	if err := decoder.Decode(&newBalance); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding request data: " + err.Error()))
		return
	}
	account, err := GOHMoneyDB.SelectAccountWithID(db, newBalance.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	originalBalance, err := account.SelectBalanceWithID(db, uint(id))
	if err == GOHMoneyDB.NoBalances {
		err = GOHMoneyDB.InvalidAccountBalanceError{
			AccountID: account.ID,
			BalanceID: uint(id),
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	updatedBalance, err := account.UpdateBalance(db, *originalBalance, newBalance.Balance)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	jsonBytes, err := json.Marshal(updatedBalance)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write(jsonBytes)
}

// accountBalanceJSONHelper is an internal type used to marshal and unmarshal json for methods.
type accountBalanceJSONHelper struct {
	AccountID uint
	Balance balance.Balance
}
