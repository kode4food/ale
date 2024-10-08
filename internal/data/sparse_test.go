package data_test

import (
	"testing"

	"github.com/kode4food/ale/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestSparseSliceSetGet(t *testing.T) {
	as := assert.New(t)
	s := data.NewSparseSlice[int]()

	s = s.Set(5, 50)
	s = s.Set(10, 100)
	s = s.Set(3, 30)

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

func TestSparseSliceSplit(t *testing.T) {
	as := assert.New(t)
	s := data.NewSparseSlice[int]()
	s = s.Set(3, 30)
	s = s.Set(7, 70)
	s = s.Set(5, 50)

	val, newSlice, ok := s.Split()
	as.True(ok)
	as.Equal(30, val)
	as.False(newSlice.Contains(3))

	_, ok = newSlice.Get(3)
	as.False(ok)
}

func TestSparseSliceClone(t *testing.T) {
	as := assert.New(t)
	s := data.NewSparseSlice[int]()
	s = s.Set(1, 10)
	s = s.Set(4, 40)
	clone := s.Clone()

	val, ok := clone.Get(1)
	as.True(ok)
	as.Equal(10, val)

	val, ok = clone.Get(4)
	as.True(ok)
	as.Equal(40, val)

	clone = clone.Set(3, 30)
	as.False(s.Contains(3))
	as.True(clone.Contains(3))
}
