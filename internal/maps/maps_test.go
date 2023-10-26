package maps_test

import (
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
