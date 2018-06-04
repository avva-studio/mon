package router

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
