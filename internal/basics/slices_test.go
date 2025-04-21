package basics_test

import (
	"cmp"
	"fmt"
	"strings"
	"testing"

	"github.com/kode4food/ale/internal/basics"
	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	as := assert.New(t)

	var s0 []string
	s1 := []string{"is", "Upper", "not", "lower"}
	s2 := []string{"is", "Upper", "not", "lower"}
	s3 := []string{"is", "Upper", "not", "lower", "extra"}
	s4 := []string{"is", "Lower", "not", "upper"}

	as.True(basics.Equal(s0, s0))
	as.False(basics.Equal(s0, s1))
	as.True(basics.Equal(s0, s1[:0]))
	as.True(basics.Equal(s1, s1))
	as.True(basics.Equal(s1, s2))
	as.False(basics.Equal(s1, s3))
	as.True(basics.Equal(s1, s3[:4]))
	as.False(basics.Equal(s1, s4))
}

func TestEqualFunc(t *testing.T) {
	as := assert.New(t)

	var s0 []string
	s1 := []string{"is", "Upper", "not", "lower"}
	s2 := []string{"is", "Upper", "not", "lower"}
	s3 := []string{"is", "Upper", "not", "lower", "extra"}
	s4 := []string{"is", "Lower", "not", "upper"}

	se := func(l, r string) bool {
		return l == r
	}

	as.True(basics.EqualFunc(s0, s0, se))
	as.False(basics.EqualFunc(s0, s1, se))
	as.True(basics.EqualFunc(s0, s1[:0], se))
	as.True(basics.EqualFunc(s1, s1, se))
	as.True(basics.EqualFunc(s1, s2, se))
	as.False(basics.EqualFunc(s1, s3, se))
	as.True(basics.EqualFunc(s1, s3[:4], se))
	as.False(basics.EqualFunc(s1, s4, se))
}

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

func TestSortedFunc(t *testing.T) {
	as := assert.New(t)
	sm := basics.SortedFunc(
		[]string{"c", "r", "b", "a"},
		func(l, r string) int {
			return -cmp.Compare(l, r)
		},
	)
	as.Equal([]string{"r", "c", "b", "a"}, sm)
}

func TestFilter(t *testing.T) {
	as := assert.New(t)
	input := []int{1, 2, 3, 5, 7, 9, 11, 12, 13, 15, 16, 17, 19}
	f1 := basics.Filter(input,
		func(i int) bool {
			return i%2 == 0
		},
	)
	as.Equal([]int{2, 12, 16}, f1)

	f2 := basics.Filter(input,
		func(i int) bool {
			return i%2 == 1
		},
	)
	as.Equal([]int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}, f2)
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

func BenchmarkEqualFunc(b *testing.B) {
	largeSlice1 := make([]int, 1_000_000)
	for i := range largeSlice1 {
		largeSlice1[i] = i
	}
	largeSlice2 := largeSlice1

	cmpFunc := func(a, b int) bool {
		return a == b
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = basics.EqualFunc(largeSlice1, largeSlice2, cmpFunc)
	}
}
