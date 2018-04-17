# GOHMoneyREST
A simple REST API for interacting with a PostgreSQL backend for tracking monetary accounts.

# BUGS
- Out of date README
- In the server, for an appHandler, the status code gets written to the header and also to the return of appHandler.ServeHttp(), which doesn't really make sense.
	- How can we make it so that the status code only gets written to a single place?