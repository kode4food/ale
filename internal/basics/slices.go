package basics

import (
	"cmp"
	"slices"
)

func Find[T any](s []T, f func(T) bool) (T, bool) {
	for _, v := range s {
		if f(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

func SortedMap[In any, Out cmp.Ordered](s []In, f func(In) Out) []Out {
	res := Map(s, f)
	slices.Sort(res)
	return res
}

func Map[In, Out any](s []In, f func(In) Out) []Out {
	res := make([]Out, len(s))
	for i, v := range s {
		res[i] = f(v)
	}
	return res
}

func IndexedMap[In, Out any](s []In, fn func(elem In, idx int) Out) []Out {
	res := make([]Out, len(s))
	for i, e := range s {
		res[i] = fn(e, i)
	}
	return res
}

func Filter[T any](s []T, f func(T) bool) []T {
	res := make([]T, 0, len(s))
	for _, v := range s {
		if f(v) {
			res = append(res, v)
		}
	}
	return slices.Clip(res)
}
