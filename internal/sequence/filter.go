package sequence

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

type FilterFunc func(ale.Value) bool

// Filter creates a new filtered Sequence
func Filter(s data.Sequence, filter FilterFunc) data.Sequence {
	var res LazyResolver
	next := s

	res = func() (ale.Value, data.Sequence, bool) {
		for f, r, ok := next.Split(); ok; f, r, ok = r.Split() {
			next = r
			if filter(f) {
				return f, NewLazy(res), true
			}
		}
		return data.Null, data.Null, false
	}
	return NewLazy(res)
}
