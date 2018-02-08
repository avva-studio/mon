package server

import (
	"net/http"
	"testing"
)

func Test_balances(t *testing.T) {
	serveFn := func(s *server, w http.ResponseWriter, r *http.Request) (int, error) {
		return s.balances(1)(w, r)
	}
	nilResponseWriterTest(t, serveFn)
	storageFuncErrorTest(t, serveFn)
}
