package time

import (
	"bytes"
	"fmt"
	"testing"
	gotime "time"

	"github.com/stretchr/testify/assert"
	"github.com/glynternet/GOHMoney/common"
)

func Test_New(t *testing.T) {
	r, err := New()
	assert.Nil(t, err)
	assert.Equal(t, NullTime{}, r.Start())
	assert.Equal(t, NullTime{}, r.End())
}

func Test_Contains(t *testing.T) {
	now := gotime.Now()
	openRange := Range{}
	openEnded := newTimeRange(t, Start(now))
	openStarted := newTimeRange(t, End(now))
	closedEnds := newTimeRange(t, Start(now), End(now.AddDate(1, 0, 0)))

	testSets := []struct {
		Range
		gotime.Time
		contains bool
	}{
		{
			Range:    openRange,
			contains: true,
		},
		{
			Range:    openRange,
			Time:     now,
			contains: true,
		},
		{
			Range:    openEnded,
			Time:     now.AddDate(-1, 0, 0),
			contains: false,
		},
		{
			Range:    openEnded,
			Time:     now,
			contains: true,
		},
		{
			Range:    openEnded,
			Time:     now.AddDate(1, 0, 0),
			contains: true,
		},
		{
			Range:    openStarted,
			Time:     now.AddDate(-1, 0, 0),
			contains: true,
		},
		{
			Range:    openStarted,
			Time:     now,
			contains: false,
		},
		{
			Range:    openStarted,
			Time:     now.AddDate(1, 0, 0),
			contains: false,
		},
		{
			Range:    closedEnds,
			Time:     now.AddDate(-2, 0, 0),
			contains: false,
		},
		{
			Range:    closedEnds,
			Time:     now,
			contains: true,
		},
		{
			Range:    closedEnds,
			Time:     now.AddDate(0, 6, 0),
			contains: true,
		},
		{
			Range:    closedEnds,
			Time:     now.AddDate(1, 0, 0),
			contains: false,
		},
		{
			Range:    closedEnds,
			Time:     now.AddDate(2, 0, 0),
			contains: false,
		},
	}
	for _, testSet := range testSets {
		contains := testSet.Range.Contains(testSet.Time)
		if contains != testSet.contains {
			var message bytes.Buffer
			fmt.Fprint(&message, `Unexpected Contains result.`)
			fmt.Fprintf(&message, "\nExpected Contains: %t\nActual Contains  : %t", testSet.contains, contains)
			fmt.Fprintf(&message, "\nTimeRange: %+v", testSet.Range)
			fmt.Fprintf(&message, "\nTime: %+v", testSet.Time)
			t.Error(message.String())
		}
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
	common.FatalIfError(t, err, "Creating timerange")
	return *r
}