package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/GlynOwenHanmer/GOHMoney"
	"github.com/GlynOwenHanmer/GOHMoneyDB"
	"github.com/gorilla/mux"
)

// Accounts handler writes a json blob for the Accounts in the DB
func Accounts(w http.ResponseWriter, r *http.Request) {
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
	accounts, err := GOHMoneyDB.SelectAccounts(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		panic(err)
	}
}

// AccountsOpen handler writes a json blob for the Accounts in the DB that are open
func AccountsOpen(w http.ResponseWriter, r *http.Request) {
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
	openAccounts, err := GOHMoneyDB.SelectAccountsOpen(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(openAccounts); err != nil {
		panic(err)
	}
}

// AccountId handler writes a json blob for the Account in the DB with the id matching the id given in the request.
// If there are any errors parsing an account ID from the request, the response code will bea 400
// If no account exists with the id, the response code will be a 404
func AccountId(w http.ResponseWriter, r *http.Request) {
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
	vars := mux.Vars(r)
	accountIdString := vars[`id`]
	if len(accountIdString) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`No id present`))
		return
	}
	accountId, err := strconv.ParseUint(accountIdString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	account, err := GOHMoneyDB.SelectAccountWithID(db, uint(accountId))
	if err != nil {
		if _, ok := err.(GOHMoneyDB.NoAccountWithIdError); ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		panic(err)
	}
}

// AccountCreate accepts a Account json blob which is then created in the backend.
// If successful, the response will contain a json blob describing the newly created Account item,
// else, the response will contain an error message describing why the creation was not successful.
func AccountCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newAccount GOHMoney.Account
	err := decoder.Decode(&newAccount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Error decoding new account: ` + err.Error()))
		return
	}
	defer r.Body.Close()
	db, err := GOHMoneyDB.OpenDBConnection(connectionString)
	if err != nil {
		ServiceUnavailableResponse(w)
		return
	}
	defer db.Close()
	createdAccount, err := GOHMoneyDB.CreateAccount(db, newAccount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`DB error creating new account: ` + err.Error()))
		return
	}
	createdAccountJson, err := json.Marshal(createdAccount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Sorry, there has been an error creating the new account.`))
		log.Println(`Error creating json from created account. New account: %s. Created account: %s`, newAccount, createdAccount)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(createdAccountJson)
}

//Retrieves all balances for an account.
func AccountBalances(w http.ResponseWriter, r *http.Request) {
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
	vars := mux.Vars(r)
	accountIdString := vars[`id`]
	if len(accountIdString) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`No id present`))
		return
	}
	accountId, err := strconv.ParseUint(accountIdString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		message := fmt.Sprintf(`Error parsing account id into uint. Error: %s`, err.Error())
		w.Write([]byte(message))
		return
	}
	account, err := GOHMoneyDB.SelectAccountWithID(db, uint(accountId))
	if _, isNoAccountError := err.(GOHMoneyDB.NoAccountWithIdError); isNoAccountError {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	balances, err := account.Balances(db)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(balances); err != nil {
		panic(err)
	}
}

// AccountBalance responds with a json blob of a Balance representing the Balance for an Account at a given date.
// The date should be url encoded.
// If no date is given, today's date is used.
func AccountBalance(w http.ResponseWriter, r *http.Request) {
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
	vars := mux.Vars(r)
	accountIdString := vars[`id`]
	if len(accountIdString) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`No id present`))
		return
	}
	accountId, err := strconv.ParseUint(accountIdString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	account, err := GOHMoneyDB.SelectAccountWithID(db, uint(accountId))
	if _, isNoAccountError := err.(GOHMoneyDB.NoAccountWithIdError); isNoAccountError {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	queryParams := r.URL.Query()
	var paramErrors []error
	var date time.Time
	for key, value := range queryParams {
		switch key {
		case `date`:
			date, err = parseDateString(value[0])
			if err != nil {
				paramErrors = append(paramErrors, err)
			}
		default:
			paramErrors = append(paramErrors, errors.New(`Invalid parameter `+key))
		}
	}
	if len(paramErrors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		var message bytes.Buffer
		fmt.Fprint(&message, `Errors: `)
		for _, err := range paramErrors {
			fmt.Fprintf(&message, `%s. `, err)
		}
		w.Write(message.Bytes())
		return
	}
	if date.IsZero() {
		date = time.Now()
	}
	balance, err := GOHMoneyDB.Account(account).BalanceAtDate(db, date)
	if err == GOHMoneyDB.NoBalances {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(balance); err != nil {
		panic(err)
	}
}

// AccountUpdate handler accepts json representing a potential update to a GOHMoneyDB.Account. The Account is decoded and attempted to be updated in the backend.
// If successful, the response contains json representing the newly updated GOHMoneyDB.Account object and returns a 204 status.
// else, an error describing why the update was unsuccessful.
func AccountUpdate(w http.ResponseWriter, r *http.Request) {
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
	decoder := json.NewDecoder(r.Body)
	var updates GOHMoney.Account
	if err := decoder.Decode(&updates); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding request data: " + err.Error()))
		return
	}
	if err := updates.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Proposed original updates are invalid: " + err.Error()))
		return
	}
	original, err := GOHMoneyDB.SelectAccountWithID(db, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	updated, err := original.Update(db, updates)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error occured updating Account: " + err.Error()))
		return
	}
	jsonBytes, err := json.Marshal(updated)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write(jsonBytes)
}
