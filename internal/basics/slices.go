package basics

import (
	"cmp"
	"slices"
)

// Equal compares two slices of the same type for equality
func Equal[T comparable](l, r []T) bool {
	if len(l) != len(r) {
		return false
	}
	if len(l) == 0 || &l[0] == &r[0] {
		return true
	}
	return slices.Equal(l, r)
}

// EqualFunc compares two slices of the same type using a custom comparison
func EqualFunc[T any](l, r []T, cmp func(T, T) bool) bool {
	if len(l) != len(r) {
		return false
	}
	if len(l) == 0 || &l[0] == &r[0] {
		return true
	}
	return slices.EqualFunc(l, r, cmp)
}

// Find returns the first element in the slice that satisfies the predicate
func Find[T any](s []T, f func(T) bool) (T, bool) {
	for _, v := range s {
		if f(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// SortedFunc returns a sorted copy of the slice using the provided comparison
// function
func SortedFunc[T any](s []T, cmp func(l, r T) int) []T {
	res := slices.Clone(s)
	slices.SortFunc(res, cmp)
	return res
}

// SortedMap returns a sorted copy of the slice after applying the provided
// mapping function to each element
func SortedMap[In any, Out cmp.Ordered](s []In, f func(In) Out) []Out {
	res := Map(s, f)
	slices.Sort(res)
	return res
}

// Map returns a new slice by applying the provided mapping function to each
// element of the input slice
func Map[In, Out any](s []In, f func(In) Out) []Out {
	res := make([]Out, len(s))
	for i, v := range s {
		res[i] = f(v)
	}
	return res
}

// IndexedMap returns a new slice by applying the provided mapping function to
// each element of the input slice, passing the index of each element as a
// second argument
func IndexedMap[In, Out any](s []In, fn func(elem In, idx int) Out) []Out {
	res := make([]Out, len(s))
	for i, e := range s {
		res[i] = fn(e, i)
	}
	return res
}

// Filter returns a new slice containing only the elements that satisfy the
// provided predicate function
func Filter[T any](s []T, f func(T) bool) []T {
	res := make([]T, 0, len(s))
	for _, v := range s {
		if f(v) {
			res = append(res, v)
		}
	}
	return slices.Clip(res)
}
