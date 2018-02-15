package server

import "net/http"

const (
	EndpointAccounts = "/accounts"
	patternAccounts  = EndpointAccounts
)

const (
	EndpointAccount       = "/account"
	EndpointFmtAccount    = EndpointAccount + "/%d"
	patternAccount        = EndpointAccount + "/{id}"
	EndpointAccountInsert = EndpointAccount + "/insert"
)

const (
	EndpointFmtAccountBalances = EndpointAccount + "/%d/balances"
	patternAccountBalances     = EndpointAccount + "/{id}/balances"
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
		{
			name:       "accountbalances",
			pattern:    patternAccountBalances,
			appHandler: s.muxAccountBalancesHandlerFunc,
			method:     http.MethodGet,
		},
	}
}
