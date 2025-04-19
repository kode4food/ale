package basics

import (
	"cmp"
	"slices"
)

func MapKeys[K comparable, V any](m map[K]V) []K {
	s := make([]K, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	s := MapKeys(m)
	slices.Sort(s)
	return s
}

func SortedKeysFunc[K comparable, V any](
	m map[K]V, cmp func(l, r K) int,
) []K {
	s := MapKeys(m)
	slices.SortFunc(s, cmp)
	return s
}
