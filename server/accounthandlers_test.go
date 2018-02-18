package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glynternet/accounting-rest/testutils"
	"github.com/glynternet/go-accounting-storagetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/glynternet/go-accounting-storage"
)

func Test_handlerSelectAccounts(t *testing.T) {
	nilResponseWriterTest(t, (*server).handlerSelectAccounts)
	storageFuncErrorTest(t, (*server).handlerSelectAccounts)

	for _, test := range []struct {
		name string
		code int
		err  error
	}{
		{
			name: "error",
			code: http.StatusServiceUnavailable,
			err:  errors.New("selecting handlerSelectAccounts"),
		},
		{
			name: "success",
			code: http.StatusOK,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			// TODO: Lots of handlers don't even require a response writer so we should remove this and force the design through the appHandlers not having a ResponseWriter. Or something like that?
			rec := httptest.NewRecorder()
			code, as, err := (&server{
				NewStorage: testutils.NewMockStorageFunc(
					&accountingtest.Storage{Err: test.err},
					false,
				),
			}).handlerSelectAccounts(rec, nil)
			assert.Equal(t, test.code, code)

			if test.err != nil {
				assert.Equal(t, test.err, errors.Cause(err))
				return
			}
			assert.NoError(t, err)
			_, ok := as.(*storage.Accounts)
			assert.True(t, ok)
		})
	}
}

func Test_handlerSelectAccount(t *testing.T) {
	serveFn := func(s *server, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
		return s.handlerSelectAccount(1)(w, r)
	}
	nilResponseWriterTest(t, serveFn)
	storageFuncErrorTest(t, serveFn)

	for _, test := range []struct {
		name string
		code int
		err  error
	}{
		{
			name: "error",
			code: http.StatusNotFound,
			err:  errors.New("selecting handlerSelectAccounts"),
		},
		{
			name: "success",
			code: http.StatusOK,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			srv := &server{
				NewStorage: testutils.NewMockStorageFunc(
					&accountingtest.Storage{AccountErr: test.err},
					false,
				),
			}
			code, a, err := serveFn(srv, rec, nil)
			assert.Equal(t, test.code, code)

			if test.err != nil {
				assert.Equal(t, test.err, errors.Cause(err))
				return
			}

			assert.NoError(t, err)
			_, ok := a.(*storage.Account)
			assert.True(t, ok)
		})
	}
}

//func Test_handlerInsertAccount(t *testing.T) {
//	serveFn := (*server).handlerInsertAccount
//	nilResponseWriterTest(t, serveFn)
//	storageFuncErrorTest(t, serveFn)
//
//	for _, test := range []struct {
//		name        string
//		body        string
//		code        int
//		errContains string
//	}{
//		{
//			name:        "unable to unmarshal into account",
//			body:        `wassssuuuuuuup?`,
//			code:        http.StatusBadRequest,
//			errContains: "unmarshalling body to account",
//		},
//	} {
//		t.Run(test.name, func(t *testing.T) {
//			w := httptest.NewRecorder()
//			r, err := http.NewRequest("any", "any", strings.NewReader(test.body))
//			common.FatalIfError(t, err, "creating new request")
//
//			var storageErr error
//			if test.errContains != "" {
//				storageErr = errors.New("")
//			}
//
//			srv := &server{
//				NewStorage: testutils.NewMockStorageFunc(
//					&accountingtest.Storage{AccountErr: storageErr},
//					false,
//				),
//			}
//
//			code, err := serveFn(srv, w, r)
//			assert.Equal(t, test.code, code)
//			if test.err {
//				assert.Contains(t, err.Error())
//			} else {
//
//			}
//		})
//	}
//}

func nilResponseWriterTest(t *testing.T, serveFunc func(*server, http.ResponseWriter, *http.Request) (int, interface{}, error)) {
	t.Run("nil response writer", func(t *testing.T) {
		code, _, err := serveFunc(&server{}, nil, nil)
		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, code)
	})
}

func storageFuncErrorTest(t *testing.T, serveFunc func(*server, http.ResponseWriter, *http.Request) (int, interface{}, error)) {
	t.Run("StorageFunc error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		code, _, err := serveFunc(
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
