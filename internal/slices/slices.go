package slices

import (
	"cmp"
	"slices"
)

func Map[In, Out any](in []In, fn func(In) Out) []Out {
	return IndexedMap(in, func(e In, _ int) Out {
		return fn(e)
	})
}

func IndexedMap[In, Out any](in []In, fn func(In, int) Out) []Out {
	res := make([]Out, len(in))
	for i, e := range in {
		res[i] = fn(e, i)
	}
	return res
}

func SortedMap[In any, Out cmp.Ordered](in []In, fn func(In) Out) []Out {
	res := Map(in, fn)
	slices.Sort(res)
	return res
}

func Filter[T any](in []T, fn func(T) bool) []T {
	res := make([]T, 0, len(in))
	for _, e := range in {
		if fn(e) {
			res = append(res, e)
		}
	}
	return slices.Clip(res)
}

func Find[T any](in []T, fn func(T) bool) (T, bool) {
	return IndexedFind(in, func(e T, _ int) bool {
		return fn(e)
	})
}

func IndexedFind[T any](in []T, fn func(T, int) bool) (T, bool) {
	for i, e := range in {
		if fn(e, i) {
			return e, true
		}
	}
	var zero T
	return zero, false
}
