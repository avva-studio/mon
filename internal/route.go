package internal

import "net/http"

const EndpointAccounts = "/accounts"

type route struct {
	pattern    string
	appHandler appHandler
}

var routes = []route{
	{
		pattern: EndpointAccounts,
		appHandler: appHandler{
			method:  http.MethodGet,
			name:    "accounts",
			handler: accounts,
		},
	},
}
