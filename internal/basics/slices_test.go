package basics_test

import (
	"cmp"
	"slices"
	"testing"

	"github.com/kode4food/ale/internal/basics"
	"github.com/stretchr/testify/assert"
)

func TestMapKeys(t *testing.T) {
	as := assert.New(t)
	k := basics.MapKeys(map[string]any{
		"age":  42,
		"name": "bob",
	})
	slices.Sort(k)
	as.Equal([]string{"age", "name"}, k)
}

func TestSortedKeys(t *testing.T) {
	as := assert.New(t)
	sk := basics.SortedKeys(map[string]any{
		"occupation": "worker bee",
		"name":       "bob",
		"age":        42,
	})
	as.Equal([]string{"age", "name", "occupation"}, sk)
}

func TestSortedKeysFunc(t *testing.T) {
	as := assert.New(t)
	sk := basics.SortedKeysFunc(map[string]any{
		"occupation": "worker bee",
		"name":       "bob",
		"age":        42,
	}, func(l, r string) int {
		return -cmp.Compare(l, r)
	})
	as.Equal([]string{"occupation", "name", "age"}, sk)
}
