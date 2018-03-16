package money_test

import (
	"testing"

	"encoding/json"

	"github.com/glynternet/go-money/currency"
	"github.com/glynternet/go-money/money"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c, err := currency.NewCode("RIN")
	assert.Nil(t, err)
	m := money.New(123, *c)
	assert.NotNil(t, m)
	assert.Equal(t, "RIN", m.Currency().String())
	assert.Equal(t, 123, m.Amount())
}

func TestJSON(t *testing.T) {
	c, err := currency.NewCode("RIN")
	assert.Nil(t, err)
	ma := money.New(9876, *c)
	bs, err := json.Marshal(ma)
	assert.Nil(t, err)
	mb, err := money.UnmarshalJSON(bs)
	assert.Nil(t, err, string(bs))
	assert.Equal(t, ma, *mb)
}

func TestUnmarshalJSON(t *testing.T) {
	i := []byte(`{"nowthen"}`)
	_, err := money.UnmarshalJSON(i)
	assert.NotNil(t, err)
	assert.IsType(t, new(json.SyntaxError), err)

	invalid := struct {
		Amount   int
		Currency string
	}{
		Amount:   12,
		Currency: "TOO_LONG",
	}
	j, err := json.Marshal(invalid)
	assert.Nil(t, err)
	_, err = money.UnmarshalJSON(j)
	assert.NotNil(t, err)
}
