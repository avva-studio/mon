package money

import (
	"encoding/json"

	"github.com/glynternet/go-money/currency"
)

// Money is an object representing a value and currency
type Money interface {
	Amount() int
	Currency() currency.Code
}

// New returns a new Money
func New(amount int, currency currency.Code) Money {
	return money{amount: amount, currency: currency}
}

type money struct {
	amount   int
	currency currency.Code
}

// Amount returns the value of the Money formed from the currency's lowest
// denominator.
// e.g. For £45.67, Amount() would return 4567
func (m money) Amount() int {
	return m.amount
}

// Currency returns the currency.Code of the Money.
func (m money) Currency() currency.Code {
	return m.currency
}

func (m money) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Amount   int
		Currency currency.Code
	}{
		Amount:   m.amount,
		Currency: m.currency,
	})
}

// UnmarshalJSON attempts to unmarshal a []byte into a money,
// returning the money, if successful, and an error, if any occurred.
func UnmarshalJSON(data []byte) (m *Money, err error) {
	var aux struct {
		Amount   int
		Currency string
	}
	err = json.Unmarshal(data, &aux)
	if err != nil {
		return nil, err
	}
	var c *currency.Code
	c, err = currency.NewCode(aux.Currency)
	if err != nil {
		return nil, err
	}
	m = new(Money)
	*m = money{
		amount:   aux.Amount,
		currency: *c,
	}
	return
}
