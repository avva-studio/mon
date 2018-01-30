package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storagetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_accounts(t *testing.T) {
	t.Run("nil response writer", func(t *testing.T) {
		code, err := (&server{}).accounts(nil, nil)
		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, code)
	})

	t.Run("StorageFunc error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		code, err := (&server{
			StorageFunc: newStorageFunc(nil, true),
		}).accounts(rec, nil)
		assert.Error(t, err)
		assert.Equal(t, testStorageFuncError, errors.Cause(err))
		assert.Equal(t, http.StatusServiceUnavailable, code)
	})

	for _, test := range []struct {
		name string
		code int
		as   *storage.Accounts
		err  error
	}{
		{
			name: "error",
			code: http.StatusServiceUnavailable,
			err:  errors.New("selecting accounts"),
		},
		{
			name: "success",
			code: http.StatusOK,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			code, err := (&server{
				StorageFunc: newStorageFunc(
					&storagetest.Storage{Accounts: test.as, Err: test.err},
					false,
				),
			}).accounts(rec, nil)
			assert.Equal(t, test.code, code)

			if test.err != nil {
				assert.Equal(t, test.err, errors.Cause(err))
				return
			}

			assert.NoError(t, err)
			ct := rec.HeaderMap[`Content-Type`]
			assert.Len(t, ct, 1)
			assert.Equal(t, `application/json; charset=UTF-8`, ct[0])
			assert.NoError(t, err)
		})
	}
}
