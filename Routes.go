package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		Name:"Index",
		Method:"GET",
		Pattern:"/",
		HandlerFunc:func(w http.ResponseWriter, r *http.Request) {w.Write([]byte(`GOHMoneyREST`))},
	},
	// Account Handlers
	Route{
		Name:"Accounts",
		Method:"GET",
		Pattern:"/accounts",
		HandlerFunc:Accounts,
	},
	Route{
		Name:        "AccountsOpen",
		Method:      "GET",
		Pattern:     "/accounts/open",
		HandlerFunc: AccountsOpen,
	},
	Route{
		Name:        "AccountId",
		Method:      "GET",
		Pattern:     "/account/{id}",
		HandlerFunc: AccountId,
	},
	Route{
		Name:        "AccountBalances",
		Method:      "GET",
		Pattern:     "/account/{id}/balances",
		HandlerFunc: AccountBalances,
	},
	Route{
		Name:        "AccountBalance",
		Method:      "GET",
		Pattern:     "/account/{id}/balance",
		HandlerFunc: AccountBalance,
	},
	Route{
		Name:        "AccountCreate",
		Method:      "POST",
		Pattern:     "/account/create",
		HandlerFunc: AccountCreate,
	},
	// Balance Handlers
	Route{
		Name:"BalanceCreate",
		Method:"POST",
		Pattern:"/balance/create",
		HandlerFunc: BalanceCreate,
	},
}
