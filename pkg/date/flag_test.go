package date

import (
	"strconv"
	"testing"
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
						assert.Nil(t, f.Time)
						assert.Error(t, err)
						return
					}
					assert.NotNil(t, f)
					diff := absDuration(f.Time.Sub(test.Time))
					acceptableThreshold := time.Millisecond * 10
					assert.Truef(t,
						diff < acceptableThreshold,
						"difference should be small but was %s. "+
							"If the difference is still quite small, "+
							"this could be because of a slow running test.",
						diff)
					assert.NoError(t, err)
				})
			}
		})
	}
}

// absDuration will make any Duration positive.
// BUG: If the duration given is equal to math.MinInt64, which is a large
// negative number, the duration returned will be negative, still.
func absDuration(a time.Duration) time.Duration {
	if a >= 0 {
		return a
	}
	return a * -1
}
