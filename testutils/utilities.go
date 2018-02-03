package testutils

import (
	"github.com/glynternet/go-accounting-storage"
	"github.com/pkg/errors"
)

func NewMockStorageFunc(s storage.Storage, err bool) func() (storage.Storage, error) {
	var rErr error
	if err {
		rErr = ErrMockStorageFunc
	}
	return func() (storage.Storage, error) {
		return s, rErr
	}
}

var ErrMockStorageFunc = errors.New("StorageFunc error")
