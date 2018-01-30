package internal

import (
	"github.com/glynternet/go-accounting-storage"
	"github.com/pkg/errors"
)

func newStorageFunc(s storage.Storage, err bool) StorageFunc {
	var rErr error
	if err {
		rErr = testStorageFuncError
	}
	return func() (storage.Storage, error) {
		return s, rErr
	}
}

var testStorageFuncError = errors.New("StorageFunc error")
