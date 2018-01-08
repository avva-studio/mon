package internal

import "net/http"

const EndpointAccounts = "/accounts"

type route struct {
	Pattern    string
	AppHandler appHandler
}

var routes = []route{
	{
		Pattern: EndpointAccounts,
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
