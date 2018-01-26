package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_accounts(t *testing.T) {
	NewStorage = mockStorage{}.storageFunc
	code, err := accounts(nil, nil)
	t.Fatal("implement this")
	assert.Nil(t, code)
	assert.Nil(t, err)
}
