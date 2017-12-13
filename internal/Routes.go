package internal

import (
	"log"
	"net/http"
)

type appHandler struct {
	name    string
	method  string
	handler func(http.ResponseWriter, *http.Request) (int, error)
}

// ServeHTTP makes our appHandler function satisfy the http.HandlerFunc interface
// If we are returning an error from our appHandler, we should not have already
// written to our ResponseWriter
func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := ah.handler(w, r); err != nil {
		log.Printf(
			"error serving on appHandler %v. Error: %v - Status: %d (%s) - Request: %+v",
			ah, err, status, http.StatusText(status), r,
		)
		switch status {
		case http.StatusServiceUnavailable:
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		// We can have cases as granular as we like, if we wanted to
		// return custom errors for specific status codes.
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			// Catch any other errors we haven't explicitly handled
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

type route struct {
	Pattern    string
	AppHandler appHandler
}

var routes = []route{
	{
		Pattern: "/accounts",
		AppHandler: appHandler{
			method:  http.MethodGet,
			name:    "Accounts",
			handler: Accounts,
		},
	},
}

//{
//	Name:        "Index",
//	Method:      "GET",
//	Pattern:     "/",
//	AppHandler: func(w http.ResponseWriter, r *http.Request) {
//		_, err := w.Write([]byte(`GOHMoneyREST`))
//		if err != nil {
//			log.Printf("Error writing response: %s", err)
//		}
//	},
//},
// Account Handlers

//{
//	Name:        "AccountsOpen",
//	Method:      "GET",
//	Pattern:     "/accounts/open",
//	AppHandler: AccountsOpen,
//},
//{
//	Name:        "AccountID",
//	Method:      "GET",
//	Pattern:     "/account/{id}",
//	AppHandler: AccountID,
//},
//{
//	Name:        "AccountBalances",
//	Method:      "GET",
//	Pattern:     "/account/{id}/balances",
//	AppHandler: AccountBalances,
//},
//{
//	Name:        "AccountBalance",
//	Method:      "GET",
//	Pattern:     "/account/{id}/balance",
//	AppHandler: AccountBalance,
//},
//{
//	Name:        "AccountCreate",
//	Method:      "POST",
//	Pattern:     "/account/create",
//	AppHandler: AccountCreate,
//},
//{
//	Name:        "AccountUpdate",
//	Method:      "PUT",
//	Pattern:     "/account/{id}/update",
//	AppHandler: AccountUpdate,
//},
//{
//	Name:        "AccountDelete",
//	Method:      "DELETE",
//	Pattern:     "/account/{id}/delete",
//	AppHandler: AccountDelete,
//},
//// Balance Handlers
//{
//	Name:        "BalanceCreate",
//	Method:      "POST",
//	Pattern:     "/balance/create",
//	AppHandler: BalanceCreate,
//},
//{
//	Name:        "BalanceUpdate",
//	Method:      "POST",
//	Pattern:     "/balance/{id}/update",
//	AppHandler: BalanceUpdate,
//},
