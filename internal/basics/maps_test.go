package basics_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kode4food/ale/internal/basics"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	as := assert.New(t)
	m := basics.Map(
		[]string{"is", "Upper", "not", "lower"},
		func(in string) bool {
			return strings.ToLower(in) != in
		},
	)
	as.Equal([]bool{false, true, false, false}, m)
}

func TestIndexedMap(t *testing.T) {
	as := assert.New(t)
	m := basics.IndexedMap(
		[]string{"is", "Upper", "not", "lower"},
		func(in string, idx int) string {
			return fmt.Sprintf("%d-%s", idx, in)
		},
	)
	as.Equal([]string{"0-is", "1-Upper", "2-not", "3-lower"}, m)
}

func TestSortedMap(t *testing.T) {
	as := assert.New(t)
	sm := basics.SortedMap([]string{"c", "r", "b", "a"},
		func(in string) string {
			return in + "-mapped"
		},
	)
	as.Equal([]string{"a-mapped", "b-mapped", "c-mapped", "r-mapped"}, sm)
}

func TestFilter(t *testing.T) {
	as := assert.New(t)
	f := basics.Filter(
		[]string{"is", "Upper", "not", "Lower"},
		func(in string) bool {
			return strings.ToLower(in) != in
		},
	)
	as.Equal([]string{"Upper", "Lower"}, f)
}

func TestFind(t *testing.T) {
	as := assert.New(t)
	f, ok := basics.Find(
		[]string{"is", "Upper", "not", "Lower"},
		func(in string) bool {
			return strings.ToLower(in) != in
		},
	)
	as.True(ok)
	as.Equal("Upper", f)
}
