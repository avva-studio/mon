package server

import "net/http"

const (
	EndpointAccounts = "/accounts"
	patternAccounts  = EndpointAccounts
)

const (
	EndpointAccount    = "/account"
	EndpointFmtAccount = EndpointAccount + "/%d"
	patternAccount     = "/account/{id}"
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
		{
			name:       "account",
			pattern:    patternAccount,
			appHandler: s.muxAccountIDHandlerFunc,
			method:     http.MethodGet,
		},
	}
}
