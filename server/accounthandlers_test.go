package server

import (
	"net/http"
	"testing"

	"github.com/glynternet/go-accounting-storagetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/glynternet/go-accounting-storage"
)

func Test_handlerSelectAccounts(t *testing.T) {

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
			server := &server{
				storage: &accountingtest.Storage{Err: test.err},
			}
			code, as, err := server.handlerSelectAccounts(nil)
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
			server := &server{
					storage: &accountingtest.Storage{AccountErr: test.err},
			}
			// request is not used in handlerSelectAccount
			code, a, err := server.handlerSelectAccount(1)(nil)
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