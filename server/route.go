package server

import "net/http"

const (
	// Accounts
	EndpointAccounts = "/accounts"
	patternAccounts  = EndpointAccounts

	// Account
	EndpointAccount       = "/account"
	EndpointFmtAccount    = EndpointAccount + "/%d"
	patternAccount        = EndpointAccount + "/{id}"
	EndpointAccountInsert = EndpointAccount + "/insert"

	// Account Balances
	EndpointFmtAccountBalances = EndpointAccount + "/%d/balances"
	patternAccountBalances     = EndpointAccount + "/{id}/balances"
)

type route struct {
	name       string
	method     string
	pattern    string
	appHandler appJSONHandler
}

func (s *server) routes() []route {
	return []route{
		{
			name:       "handlerSelectAccounts",
			pattern:    patternAccounts,
			appHandler: s.handlerSelectAccounts,
			method:     http.MethodGet,
		},
		{
			name:       "handlerSelectAccount",
			pattern:    patternAccount,
			appHandler: s.muxAccountIDHandlerFunc,
			method:     http.MethodGet,
		},
		{
			name:       "accountBalances",
			pattern:    patternAccountBalances,
			appHandler: s.muxAccountBalancesHandlerFunc,
			method:     http.MethodGet,
		},
		//{
		//	name:       "insertAccount",
		//	pattern:    EndpointAccountInsert,
		//	appHandler: s.handlerInsertAccount,
		//	method:     http.MethodPost,
		//},
	}
}
