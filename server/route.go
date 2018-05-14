package server

import "net/http"

const (
	// Accounts
	EndpointAccounts = "/accounts"
	patternAccounts  = EndpointAccounts

	// Account
	EndpointAccount          = "/account"
	EndpointFmtAccount       = EndpointAccount + "/%d"
	patternAccount           = EndpointAccount + "/{id}"
	EndpointAccountInsert    = EndpointAccount + "/insert"
	EndpointFmtAccountUpdate = EndpointFmtAccount + "/update"
	patternAccountUpdate     = patternAccount + "/update"

	// Account Balances
	EndpointFmtAccountBalances      = EndpointAccount + "/%d/balances"
	patternAccountBalances          = EndpointAccount + "/{id}/balances"
	EndpointFmtAccountBalanceInsert = EndpointAccount + "/%d/balance/insert"
	patternAccountBalanceInsert     = EndpointAccount + "/{id}/balance/insert"
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
			name:       "Accounts",
			pattern:    patternAccounts,
			appHandler: s.handlerSelectAccounts,
			method:     http.MethodGet,
		},
		{
			name:       "Account",
			pattern:    patternAccount,
			appHandler: s.muxAccountIDHandlerFunc,
			method:     http.MethodGet,
		},
		{
			name:       "AccountInsert",
			pattern:    EndpointAccountInsert,
			appHandler: s.muxAccountInsertHandlerFunc,
			method:     http.MethodPost,
		},
		{
			name:       "AccountUpdate",
			pattern:    patternAccountUpdate,
			appHandler: s.muxAccountUpdateHandlerFunc,
			method:     http.MethodPost,
		},
		{
			name:       "Balances",
			pattern:    patternAccountBalances,
			appHandler: s.muxAccountBalancesHandlerFunc,
			method:     http.MethodGet,
		},
		{
			name:       "BalanceInsert",
			pattern:    patternAccountBalanceInsert,
			appHandler: s.muxAccountBalanceInsertHandlerFunc,
			method:     http.MethodPost,
		},
	}
}
