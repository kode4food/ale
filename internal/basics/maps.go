package basics

import (
	"cmp"
	"slices"
)

// MapKeys returns the keys of the provided map as a slice
func MapKeys[K comparable, V any](m map[K]V) []K {
	s := make([]K, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

// SortedKeys returns the keys of the provided map as a sorted slice
func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	s := MapKeys(m)
	slices.Sort(s)
	return s
}

// SortedKeysFunc returns the keys of the provided map as a sorted slice
// using the provided comparison function
func SortedKeysFunc[K comparable, V any](
	m map[K]V, cmp func(l, r K) int,
) []K {
	s := MapKeys(m)
	slices.SortFunc(s, cmp)
	return s
}
