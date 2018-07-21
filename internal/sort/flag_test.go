package sort_test

import (
	"testing"

	"github.com/glynternet/mon/internal/sort"
	"github.com/stretchr/testify/assert"
)

func TestFlag_Set(t *testing.T) {
	for _, test := range []struct {
		name string
		val  string
		err  bool
		key  string
	}{
		{
			name: "zero-values",
		},
		{
			name: "valid",
			val:  "id",
			key:  "id",
		},
		{
			name: "invalid",
			val:  "alishdgkajhsdga",
			err:  true,
			key:  "id",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			k := sort.NewKey()
			err := k.Set(test.val)
			if test.err {
				assert.Error(t, err)
				assert.Equal(t, "", k.String())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.key, k.String())
		})
	}
}
