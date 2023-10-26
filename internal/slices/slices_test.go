package slices_test

import (
	"strings"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/slices"
)

func TestMap(t *testing.T) {
	as := assert.New(t)
	m := slices.Map(
		[]string{"is", "Upper", "not", "lower"},
		func(in string) bool {
			return strings.ToLower(in) != in
		},
	)
	as.Equal([]bool{false, true, false, false}, m)
}

func TestSortedMap(t *testing.T) {
	as := assert.New(t)

	m := slices.SortedMap(
		[]string{"c", "r", "b", "a"},
		func(in string) string {
			return in + "-mapped"
		},
	)
	as.Equal([]string{"a-mapped", "b-mapped", "c-mapped", "r-mapped"}, m)
}

func TestFilter(t *testing.T) {
	as := assert.New(t)
	f := slices.Filter(
		[]string{"is", "Upper", "not", "Lower"},
		func(in string) bool {
			return strings.ToLower(in) != in
		},
	)
	as.Equal([]string{"Upper", "Lower"}, f)
}

func TestFind(t *testing.T) {
	as := assert.New(t)
	s, ok := slices.Find(
		[]string{"is", "Upper", "not", "Lower"},
		func(in string) bool {
			return strings.ToLower(in) != in
		},
	)
	as.True(ok)
	as.Equal("Upper", s)
}
