package internal

import (
	"iter"
	"math/bits"
	"slices"
)

// SparseSlice manages a sparse slice of elements, only allocating enough space
// to satisfy the current set. It uses a 64-bit int as its mask, and is limited
// to 64 elements. It performs no range checking, so use wisely
type SparseSlice[T any] struct {
	data []T
	mask uint64
}

// NewSparseSlice initializes an empty SparseSlice
func NewSparseSlice[T any]() *SparseSlice[T] {
	return nil
}

// Set returns a new SparseSlice with the value set at the specified index.
func (s *SparseSlice[T]) Set(idx int, value T) *SparseSlice[T] {
	if !s.Contains(idx) {
		return s.insert(idx, value)
	}
	return s.replace(idx, value)
}

func (s *SparseSlice[T]) insert(idx int, value T) *SparseSlice[T] {
	if s == nil {
		return &SparseSlice[T]{
			data: []T{value},
			mask: 1 << idx,
		}
	}
	pos := bits.OnesCount64(s.mask & ((1 << idx) - 1))
	return &SparseSlice[T]{
		data: slices.Insert(s.data, pos, value),
		mask: s.mask | (1 << idx),
	}
}

func (s *SparseSlice[T]) replace(idx int, value T) *SparseSlice[T] {
	data := slices.Clone(s.data)
	pos := s.position(idx)
	data[pos] = value
	return &SparseSlice[T]{
		data: data,
		mask: s.mask,
	}
}

// Get retrieves a value at a specific index, returning false if it’s not set.
func (s *SparseSlice[T]) Get(idx int) (T, bool) {
	if !s.Contains(idx) {
		var zero T
		return zero, false
	}
	pos := s.position(idx)
	return s.data[pos], true
}

// Unset returns a new SparseSlice with the specified index possibly removed.
func (s *SparseSlice[T]) Unset(idx int) *SparseSlice[T] {
	if !s.Contains(idx) {
		return s
	}
	mask := s.mask & ^(1 << idx)
	if mask == 0 {
		return nil
	}
	pos := s.position(idx)
	return &SparseSlice[T]{
		data: slices.Concat(s.data[:pos], s.data[pos+1:]),
		mask: mask,
	}
}

func (s *SparseSlice[T]) IsEmpty() bool {
	return s == nil || s.mask == 0
}

func (s *SparseSlice[T]) HighIndex() int {
	if s == nil {
		return -1
	}
	return 63 - bits.LeadingZeros64(s.mask)
}

func (s *SparseSlice[T]) LowIndex() int {
	if s == nil {
		return -1
	}
	return bits.TrailingZeros64(s.mask)
}

func (s *SparseSlice[T]) Contains(idx int) bool {
	return s != nil && (s.mask&(1<<idx)) != 0
}

func (s *SparseSlice[T]) RawData() ([]T, uint64) {
	if s == nil {
		return nil, 0
	}
	return s.data, s.mask
}

func (s *SparseSlice[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		if s == nil {
			return
		}
		low, high := s.LowIndex(), s.HighIndex()
		for i := low; i <= high; i++ {
			if s.mask&(1<<i) != 0 {
				if !yield(i, s.data[s.position(i)]) {
					return
				}
			}
		}
	}
}

func (s *SparseSlice[T]) Values() iter.Seq[T] {
	if s == nil {
		return slices.Values([]T(nil))
	}
	return slices.Values(s.data)
}

func (s *SparseSlice[T]) position(idx int) int {
	return bits.OnesCount64(s.mask & ((1 << idx) - 1))
}
