package time

import (
	"testing"
	gtime "time"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	r, err := New()
	assert.Nil(t, err)
	assert.Equal(t, NullTime{}, r.Start())
	assert.Equal(t, NullTime{}, r.End())
}

func TestRange_Contains(t *testing.T) {
	now := gtime.Now()
	var openRange Range
	openEnded := newTimeRange(t, Start(now))
	openStarted := newTimeRange(t, End(now))
	closedEnds := newTimeRange(t, Start(now), End(now.AddDate(1, 0, 0)))

	for _, test := range []struct {
		name string
		Range
		gtime.Time
		contains bool
	}{
		{
			name:     "zero-values",
			contains: true,
		},
		{
			name:     "fully open any time",
			Range:    openRange,
			Time:     now,
			contains: true,
		},
		{
			name:     "open ended before start",
			Range:    openEnded,
			Time:     now.AddDate(-1, 0, 0),
			contains: false,
		},
		{
			name:     "open ended on start",
			Range:    openEnded,
			Time:     now,
			contains: true,
		},
		{
			name:     "open ended after start",
			Range:    openEnded,
			Time:     now.AddDate(1, 0, 0),
			contains: true,
		},
		{
			name:     "open started before end",
			Range:    openStarted,
			Time:     now.AddDate(-1, 0, 0),
			contains: true,
		},
		{
			name:     "open started on end",
			Range:    openStarted,
			Time:     now,
			contains: true,
		},
		{
			name:     "open started after end",
			Range:    openStarted,
			Time:     now.AddDate(1, 0, 0),
			contains: false,
		},
		{
			name:     "closed ends before start",
			Range:    closedEnds,
			Time:     now.AddDate(-2, 0, 0),
			contains: false,
		},
		{
			name:     "closed ends on start",
			Range:    closedEnds,
			Time:     now,
			contains: true,
		},
		{
			name:     "closed ends in middle",
			Range:    closedEnds,
			Time:     now.AddDate(0, 6, 0),
			contains: true,
		},
		{
			name:     "closed ends on end",
			Range:    closedEnds,
			Time:     now.AddDate(1, 0, 0),
			contains: true,
		},
		{
			name:     "closed ends after end",
			Range:    closedEnds,
			Time:     now.AddDate(2, 0, 0),
			contains: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t,
				test.contains,
				test.Range.Contains(test.Time),
				"Range: %+v\nTime: %+v", test.Range, test.Time,
			)
		})
	}
}

func Test_Equal(t *testing.T) {
	testSets := []struct {
		a, b  Range
		equal bool
	}{
		{
			a:     Range{},
			b:     Range{},
			equal: true,
		},
		{
			a: Range{
				start: NullTime{
					Valid: true,
				},
			},
			b:     Range{},
			equal: false,
		},
		{
			a: Range{
				end: NullTime{
					Valid: true,
				},
			},
			b:     Range{},
			equal: false,
		},
	}
	for _, testSet := range testSets {
		if equal := testSet.a.Equal(testSet.b); equal != testSet.equal {
			t.Errorf(`Unexpected Equal result.\nExpected: %t, Actual  : %t`, testSet.equal, equal)
		}
	}
}

func newTimeRange(t *testing.T, os ...Option) Range {
	r, err := New(os...)
	if err == nil {
		return *r
	}
	t.Fatalf("creating timerange: %s", err)
	return Range{}
}
