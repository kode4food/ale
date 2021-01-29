package sequence

import "github.com/kode4food/ale/data"

// Filter creates a new filtered Sequence
func Filter(s data.Sequence, filter data.Function) data.Sequence {
	var res LazyResolver
	next := s

	res = func() (data.Value, data.Sequence, bool) {
		for f, r, ok := next.Split(); ok; f, r, ok = r.Split() {
			next = r
			if data.Truthy(filter.Call(f)) {
				return f, NewLazy(res), true
			}
		}
		return data.Nil, data.EmptyList, false
	}
	return NewLazy(res)
}
