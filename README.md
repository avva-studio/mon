# GOHMoneyREST
A simple REST API for interacting with a PostgreSQL backend for tracking monetary accounts.

### BUGS
- When using the flag --open, the accounts are always filtered by the ones that are open currently. This is because of design in the `go-accounting` package. 