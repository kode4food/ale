package maps

func Keys[K comparable, V any](in map[K]V) []K {
	res := make([]K, len(in))
	i := 0
	for k := range in {
		res[i] = k
		i++
	}
	return res
}
