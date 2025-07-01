package sequence

import "github.com/kode4food/ale/pkg/data"

// Last returns the final element of a Sequence, possibly by scanning
func Last(s data.Sequence) (data.Value, bool) {
	if s.IsEmpty() {
		return data.Null, false
	}

	if i, ok := s.(data.Indexed); ok {
		return i.ElementAt(i.Count() - 1)
	}

	var res data.Value
	var lok bool
	for f, s, ok := s.Split(); ok; f, s, ok = s.Split() {
		res = f
		lok = ok
	}
	return res, lok
}
