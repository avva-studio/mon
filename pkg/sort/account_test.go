package sort_test

import (
	"testing"

	"github.com/glynternet/mon/pkg/sort"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	t.Run("nil input remains nil", func(t *testing.T) {
		var as storage.Accounts
		sort.ID(as)
		assert.Nil(t, as)
	})

	for _, test := range []struct {
		name    string
		in, out []uint
	}{
		{
			name: "empty storage.Accounts",
			in:   []uint{},
		},
		{
			name: "single account yields same account",
			in:   []uint{1},
			out:  []uint{1},
		},
		{
			name: "in-order input yields same order output",
			in:   []uint{0, 1, 2, 3},
			out:  []uint{0, 1, 2, 3},
		},
		{
			name: "in reverse order",
			in:   []uint{3, 2, 1, 0},
			out:  []uint{0, 1, 2, 3},
		},
		{
			name: "mixed order",
			in:   []uint{3, 7, 99, 2, 5, 8},
			out:  []uint{2, 3, 5, 7, 8, 99},
		},
		{
			name: "mixed order with repeating",
			in:   []uint{9, 5, 2, 66, 4, 5},
			out:  []uint{2, 4, 5, 5, 9, 66},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var as storage.Accounts
			for _, id := range test.in {
				as = append(as, storage.Account{ID: id})
			}

			sort.ID(as)

			var out []uint
			for _, a := range as {
				out = append(out, a.ID)
			}

			assert.Equal(t, test.out, out)
		})
	}
}
