package server

import "net/http"

const (
	EndpointAccounts = "/accounts"
	patternAccounts  = EndpointAccounts
)

const (
	EndpointAccount = "/account"
	patternAccount  = "/account/{id}"
)

type route struct {
	name       string
	method     string
	pattern    string
	appHandler appHandler
}

func (s *server) routes() []route {
	return []route{
		{
			name:       "accounts",
			pattern:    patternAccounts,
			appHandler: s.accounts,
			method:     http.MethodGet,
		},
		//{
		//	name:       "account",
		//	pattern:    patternAccount,
		//	appHandler: s.muxAccountIDHandlerfunc,
		//	method:     http.MethodGet,
		//},
	}
}
