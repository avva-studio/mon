package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glynternet/accounting-rest/testutils"
	"github.com/glynternet/go-accounting-storagetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_accounts(t *testing.T) {
	nilResponseWriterTest(t, (*server).accounts)
	storageFuncErrorTest(t, (*server).accounts)

	for _, test := range []struct {
		name string
		code int
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
				NewStorage: testutils.NewMockStorageFunc(
					&accountingtest.Storage{Err: test.err},
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

func Test_accountHandlerWithID(t *testing.T) {
	c := func(s *server, w http.ResponseWriter, r *http.Request) (int, error) {
		return s.accountHandlerWithID(1)(w, r)
	}
	nilResponseWriterTest(t, c)
	storageFuncErrorTest(t, c)

	for _, test := range []struct {
		name string
		code int
		err  error
	}{
		{
			name: "error",
			code: http.StatusServiceUnavailable,
			err:  errors.New("selecting accounts"),
		},
		//{
		//	name: "success",
		//	code: http.StatusOK,
		//},
	} {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			code, err := (&server{
				NewStorage: testutils.NewMockStorageFunc(
					&accountingtest.Storage{Err: test.err},
					false,
				),
			}).accountHandlerWithID(1)(rec, nil)
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

func nilResponseWriterTest(t *testing.T, serveFunc func(*server, http.ResponseWriter, *http.Request) (int, error)) {
	t.Run("nil response writer", func(t *testing.T) {
		code, err := serveFunc(&server{}, nil, nil)
		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, code)
	})
}

func storageFuncErrorTest(t *testing.T, serveFunc func(*server, http.ResponseWriter, *http.Request) (int, error)) {
	t.Run("StorageFunc error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		code, err := serveFunc(
			&server{
				NewStorage: testutils.NewMockStorageFunc(nil, true),
			},
			rec,
			nil,
		)
		assert.Error(t, err)
		assert.Equal(t, testutils.ErrMockStorageFunc, errors.Cause(err))
		assert.Equal(t, http.StatusServiceUnavailable, code)
	})
}
