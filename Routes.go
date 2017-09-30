package main

import (
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []route{
	{
		Name:        "Index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`GOHMoneyREST`)) },
	},
	// Account Handlers
	{
		Name:        "Accounts",
		Method:      "GET",
		Pattern:     "/accounts",
		HandlerFunc: Accounts,
	},
	{
		Name:        "AccountsOpen",
		Method:      "GET",
		Pattern:     "/accounts/open",
		HandlerFunc: AccountsOpen,
	},
	{
		Name:        "AccountID",
		Method:      "GET",
		Pattern:     "/account/{id}",
		HandlerFunc: AccountID,
	},
	{
		Name:        "AccountBalances",
		Method:      "GET",
		Pattern:     "/account/{id}/balances",
		HandlerFunc: AccountBalances,
	},
	{
		Name:        "AccountBalance",
		Method:      "GET",
		Pattern:     "/account/{id}/balance",
		HandlerFunc: AccountBalance,
	},
	{
		Name:        "AccountCreate",
		Method:      "POST",
		Pattern:     "/account/create",
		HandlerFunc: AccountCreate,
	},
	{
		Name:        "AccountUpdate",
		Method:      "PUT",
		Pattern:     "/account/{id}/update",
		HandlerFunc: AccountUpdate,
	},
	{
		Name:        "AccountDelete",
		Method:      "DELETE",
		Pattern:     "/account/{id}/delete",
		HandlerFunc: AccountDelete,
	},
	// Balance Handlers
	{
		Name:        "BalanceCreate",
		Method:      "POST",
		Pattern:     "/balance/create",
		HandlerFunc: BalanceCreate,
	},
	{
		Name:        "BalanceUpdate",
		Method:      "POST",
		Pattern:     "/balance/{id}/update",
		HandlerFunc: BalanceUpdate,
	},
}
