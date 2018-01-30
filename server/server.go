package server

import (
	"net/http"

	"github.com/pkg/errors"
)

func New(sf StorageFunc) (*server, error) {
	if sf == nil {
		return nil, errors.New("nil StorageFunc provided")
	}
	return &server{NewStorage: sf}, nil
}

type server struct {
	NewStorage StorageFunc
}

func (s *server) ListenAndServe(addr string) error {
	//logDBState()  //TODO: logDBState?
	router, err := s.newRouter()
	if err != nil {
		return errors.Wrap(err, "creating new Router")
	}
	return http.ListenAndServe(addr, router)
}
