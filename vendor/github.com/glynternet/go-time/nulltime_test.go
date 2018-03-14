package time

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestNullTime_Equal(t *testing.T) {
	timeNow := time.Now()
	testSets := []struct {
		a, b  NullTime
		equal bool
	}{
		// a not Valid, b not Valid
		{
			a:     NullTime{},
			b:     NullTime{},
			equal: true,
		},
		{
			a: NullTime{
				Time: timeNow,
			},
			b:     NullTime{},
			equal: false,
		},
		{
			a: NullTime{},
			b: NullTime{
				Time: timeNow,
			},
			equal: false,
		},
		{
			a: NullTime{
				Time: timeNow,
			},
			b: NullTime{
				Time: timeNow,
			},
			equal: true,
		},

		// a Valid, b not Valid
		{
			a: NullTime{
				Valid: true,
			},
			b:     NullTime{},
			equal: false,
		},
		{
			a: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			b:     NullTime{},
			equal: false,
		},
		{
			a: NullTime{
				Valid: true,
			},
			b: NullTime{
				Time: timeNow,
			},
			equal: false,
		},
		{
			a: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			b: NullTime{
				Time: timeNow,
			},
			equal: false,
		},

		// a not Valid, b Valid
		{
			a: NullTime{},
			b: NullTime{
				Valid: true,
			},
			equal: false,
		},
		{
			a: NullTime{
				Time: timeNow,
			},
			b: NullTime{
				Valid: true,
			},
			equal: false,
		},
		{
			a: NullTime{},
			b: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			equal: false,
		},
		{
			a: NullTime{
				Time: timeNow,
			},
			b: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			equal: false,
		},

		// a Valid, b Valid
		{
			a: NullTime{
				Valid: true,
			},
			b: NullTime{
				Valid: true,
			},
			equal: true,
		},
		{
			a: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			b: NullTime{
				Valid: true,
			},
			equal: false,
		},
		{
			a: NullTime{
				Valid: true,
			},
			b: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			equal: false,
		},
		{
			a: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			b: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			equal: true,
		},
	}

	for _, ts := range testSets {
		if equal := ts.a.Equal(ts.b); ts.equal != equal {
			var message bytes.Buffer
			fmt.Fprintf(&message, "Unexpected equal result.\nExpected: %t\nActual  : %t", ts.equal, equal)
			fmt.Fprintf(&message, "\na: %+v", ts.a)
			fmt.Fprintf(&message, "\nb: %+v", ts.b)
			t.Errorf(message.String())
		}
	}
}

func TestNullTime_EqualTime(t *testing.T) {
	now := time.Now()
	testSets := []struct {
		a  NullTime
		bs map[time.Time]bool
	}{
		{
			a: NullTime{
				Valid: false,
				Time:  now,
			},
			bs: map[time.Time]bool{
				now: false,
			},
		},
		{
			a: NullTime{
				Valid: true,
				Time:  now,
			},
			bs: map[time.Time]bool{
				now.Add(-1 * time.Millisecond): false,
				now: true,
				now.Add(1 * time.Millisecond): false,
			},
		},
	}
	for _, ts := range testSets {
		for b, e := range ts.bs {
			if equal := ts.a.EqualTime(b); equal != e {
				var message bytes.Buffer
				fmt.Fprintf(&message, "Unexpected equal result.\nExpected: %t\nActual  : %t", e, equal)
				fmt.Fprintf(&message, "\na: %+v", ts.a)
				fmt.Fprintf(&message, "\nb: %+v", b)
				t.Errorf(message.String())
			}
		}
	}
}
