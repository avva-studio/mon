package server

import "net/http"

const EndpointAccounts = "/accounts"

type route struct {
	pattern    string
	appHandler appHandler
}

func (s *server) routes() []route {
	return []route{
		{
			pattern: EndpointAccounts,
			appHandler: appHandler{
				method:  http.MethodGet,
				name:    "accounts",
				handler: s.accounts,
			},
		},
	}
}
