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
