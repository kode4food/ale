package data_test

import (
	"testing"

	"github.com/kode4food/ale/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestSparseSliceSetGet(t *testing.T) {
	as := assert.New(t)
	s := data.NewSparseSlice[int]()
	as.True(s.IsEmpty())

	s = s.Set(5, 50)
	s = s.Set(10, 100)
	s = s.Set(3, 30)
	as.False(s.IsEmpty())

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
	s = s.Set(2, 20)
	s = s.Set(6, 60)

	s = s.Unset(2)
	as.False(s.Contains(2))

	val, ok := s.Get(2)
	as.False(ok)
	as.Equal(0, val)
}

func TestSparseSliceReplace(t *testing.T) {
	as := assert.New(t)
	s := data.NewSparseSlice[int]()
	s = s.Set(3, 30)
	s = s.Set(5, 50)

	s = s.Set(3, 300)

	val, ok := s.Get(3)
	as.True(ok)
	as.Equal(300, val)

	val, ok = s.Get(5)
	as.True(ok)
	as.Equal(50, val)
}
