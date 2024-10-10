package data_test

import (
	"testing"

	"github.com/kode4food/ale/internal/data"
	"github.com/stretchr/testify/assert"
)

func testDataAll[T any](t *testing.T, s *data.SparseSlice[T], m map[int]T) {
	as := assert.New(t)
	l := 0
	for idx, v := range s.All() {
		r, ok := m[idx]
		as.True(ok)
		as.Equal(r, v)
		l++
	}
	as.Equal(l, len(m))
}

func testDataValues[T any](t *testing.T, s *data.SparseSlice[T], d []T) {
	as := assert.New(t)
	i := 0
	for v := range s.Values() {
		as.Equal(d[i], v)
		i++
	}
	as.Equal(i, len(d))
}

func TestEmptySparseSlice(t *testing.T) {
	as := assert.New(t)
	s := data.NewSparseSlice[int]()
	as.True(s.IsEmpty())
	testDataValues(t, s, []int{})
	as.Equal(-1, s.LowIndex())
	as.Equal(-1, s.HighIndex())
}

func TestSparseSliceSetGet(t *testing.T) {
	as := assert.New(t)
	s := data.NewSparseSlice[int]()
	as.True(s.IsEmpty())

	s = s.Set(5, 50)
	s = s.Set(10, 100)
	s = s.Set(3, 30)
	as.False(s.IsEmpty())

	testDataAll(t, s, map[int]int{3: 30, 5: 50, 10: 100})
	testDataValues(t, s, []int{30, 50, 100})
	as.Equal(3, s.LowIndex())
	as.Equal(10, s.HighIndex())

	val, ok := s.Get(5)
	as.True(ok)
	as.Equal(50, val)

	val, ok = s.Get(10)
	as.True(ok)
	as.Equal(100, val)

	val, ok = s.Get(3)
	as.True(ok)
	as.Equal(30, val)

	_, ok = s.Get(1)
	as.False(ok)
}

func TestSparseSliceUnset(t *testing.T) {
	as := assert.New(t)
	s := data.NewSparseSlice[int]()
	s = s.Unset(32)
	as.True(s.IsEmpty())

	s = s.Set(2, 20)
	s = s.Set(6, 60)

	testDataAll(t, s, map[int]int{2: 20, 6: 60})
	testDataValues(t, s, []int{20, 60})
	as.Equal(2, s.LowIndex())
	as.Equal(6, s.HighIndex())

	s = s.Unset(2)
	as.False(s.Contains(2))

	val, ok := s.Get(2)
	as.False(ok)
	as.Equal(0, val)

	testDataAll(t, s, map[int]int{6: 60})
	testDataValues(t, s, []int{60})
	as.Equal(6, s.LowIndex())
	as.Equal(6, s.HighIndex())

	s1 := s.Unset(10)
	as.Equal(s, s1)

	s = s.Unset(6)
	as.True(s.IsEmpty())
	as.Nil(s)
}

func TestSparseSliceReplace(t *testing.T) {
	as := assert.New(t)
	s := data.NewSparseSlice[int]()
	s = s.Set(3, 30)
	s = s.Set(5, 50)
	testDataAll(t, s, map[int]int{3: 30, 5: 50})
	testDataValues(t, s, []int{30, 50})
	s = s.Set(3, 300)

	val, ok := s.Get(3)
	as.True(ok)
	as.Equal(300, val)

	val, ok = s.Get(5)
	as.True(ok)
	as.Equal(50, val)

	testDataAll(t, s, map[int]int{3: 300, 5: 50})
	testDataValues(t, s, []int{300, 50})
}
