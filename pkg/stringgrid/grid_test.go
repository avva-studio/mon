package stringgrid

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	cs := Columns{
		SimpleColumn(func(i uint) string {
			return strconv.Itoa(int(i))
		}),
		SimpleColumn(func(i uint) string {
			return strconv.Itoa(int(i * 2))
		}),
	}

	data, err := cs.Generate(10)
	assert.NoError(t, err)
	t.Log(data)
}
