package maps_test

import (
	"cmp"
	"slices"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/maps"
)

func TestKeys(t *testing.T) {
	as := assert.New(t)
	m := maps.Keys(map[string]any{
		"age":  42,
		"name": "bob",
	})
	slices.Sort(m)
	as.Equal([]string{"age", "name"}, m)
}

func TestSortedKeys(t *testing.T) {
	as := assert.New(t)
	m := maps.SortedKeys(map[string]any{
		"occupation": "worker bee",
		"name":       "bob",
		"age":        42,
	})
	as.Equal([]string{"age", "name", "occupation"}, m)
}

func TestSortedKeysFunc(t *testing.T) {
	as := assert.New(t)
	m := maps.SortedKeysFunc(map[string]any{
		"occupation": "worker bee",
		"name":       "bob",
		"age":        42,
	}, func(l, r string) int {
		return -cmp.Compare(l, r)
	})
	as.Equal([]string{"occupation", "name", "age"}, m)
}
