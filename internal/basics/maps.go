package basics

import (
	"cmp"
	"slices"
)

// MapKeys returns the keys of the provided map as a slice
func MapKeys[K comparable, V any](m map[K]V) []K {
	s := make([]K, len(m))
	var i int
	for k := range m {
		s[i] = k
		i++
	}
	return s
}

// MapValues returns the values of the provided map as a slice
func MapValues[K comparable, V any](m map[K]V) []V {
	s := make([]V, len(m))
	var i int
	for _, v := range m {
		s[i] = v
		i++
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
func SortedKeysFunc[K comparable, V any](m map[K]V, cmp func(l, r K) int) []K {
	s := MapKeys(m)
	slices.SortFunc(s, cmp)
	return s
}
