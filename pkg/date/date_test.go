package date

import (
	"testing"

	"strconv"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestFlag_Set(t *testing.T) {
	for _, test := range []struct {
		name string
		vals []string
		err  bool
		time.Time
	}{
		{
			name: "zero-values",
			err:  true,
		},
		{
			name: "whitespace",
			vals: []string{"\t\t\t\r"},
			err:  true,
		},
		{
			name: "yesterday",
			vals: []string{"yesterday", "YESTERDAY", "y", "Y"},
			Time: time.Now().Add(time.Hour * -24),
		},
		{
			name: "nonsense",
			vals: []string{"bloopy bleep", "!!!!!"},
			err:  true,
		},
		{
			name: "invalid date format",
			vals: []string{"02/03", "02-03", "-1000-01-87"},
			err:  true,
		},
		{
			name: "valid",
			vals: []string{"2018-03-02"},
			Time: time.Date(2018, 03, 02, 0, 0, 0, 0, time.UTC),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			for key, val := range test.vals {
				k := key
				v := val
				t.Run(strconv.Itoa(k), func(t *testing.T) {
					f := &flag{}
					err := f.Set(v)
					if test.err {
						assert.Error(t, err)
						return
					}
					assert.NotNil(t, f)
					diff := absDuration(time.Time(*f).Sub(test.Time))
					assert.Truef(t, diff < time.Millisecond*5 && diff > -time.Millisecond*5, "difference should be small but was %d", diff)
					assert.NoError(t, err)
				})
			}
		})
	}
}

func absDuration(a time.Duration) time.Duration {
	if a >= 0 {
		return a
	}
	return a * -1
}
