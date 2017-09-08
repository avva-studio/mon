package main

import (
	"net/http"
	"github.com/GlynOwenHanmer/GOHMoney"
	"github.com/GlynOwenHanmer/GOHMoneyDB"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
)

// BalanceCreate handler accepts json representing a potential new GOHMoney.Balance. The Balance is decoded and attempted to be added to the backend.
// If successful, the response contains json representing the newly created GOHMoneyDB.Balance object,
// else, an error describing why the creation was unsuccessful.
func BalanceCreate(w http.ResponseWriter, r *http.Request)  {
	type accountBalance struct {
		AccountId uint `json:"account_id"`
		GOHMoney.Balance `json:"balance"`
	}
	decoder := json.NewDecoder(r.Body)
	var newBalance accountBalance
	err := decoder.Decode(&newBalance)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding new account: " + err.Error()))
		return
	}
	defer r.Body.Close()
	if newBalance.AccountId < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Account ID must be a positive integer`))
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
	account, err := GOHMoneyDB.SelectAccountWithID(db, newBalance.AccountId)
	if err != nil {
		if _, ok := err.(GOHMoneyDB.NoAccountWithIdError); ok {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}
	if err := newBalance.Validate(); err != nil {
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
	id, err := parseIdString(vars[`id`])
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
	type accountBalance struct {
		AccountId uint
		GOHMoney.Balance
	}
	decoder := json.NewDecoder(r.Body)
	var newBalance accountBalance
	if err := decoder.Decode(&newBalance); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding request data: " + err.Error()))
		return
	}
	account, err := GOHMoneyDB.SelectAccountWithID(db, newBalance.AccountId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	originalBalance, err := account.SelectBalanceWithId(db,uint(id))
	if err == GOHMoneyDB.NoBalances {
		err = GOHMoneyDB.InvalidAccountBalanceError{
			AccountId:account.Id,
			BalanceId:uint(id),
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