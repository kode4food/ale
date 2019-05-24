package stdlib

import "gitlab.com/kode4food/ale/data"

// Filter creates a new filtered Sequence
func Filter(s data.Sequence, filter data.Call) data.Sequence {
	var res LazyResolver
	next := s

	res = func() (data.Value, data.Sequence, bool) {
		for f, r, ok := next.Split(); ok; f, r, ok = r.Split() {
			next = r
			if data.Truthy(filter(f)) {
				return f, NewLazySequence(res), true
			}
		}
		return data.Nil, data.EmptyList, false
	}
	return NewLazySequence(res)
}

// Concat creates a new sequence based on the content of several Sequences
func Concat(s ...data.Value) data.Sequence {
	var res LazyResolver
	var next data.Sequence = data.NewVector(s...)

	res = func() (data.Value, data.Sequence, bool) {
		for f, r, ok := next.Split(); ok; f, r, ok = r.Split() {
			v := f.(data.Sequence)
			next = r
			if vf, vr, ok := v.Split(); ok {
				next = next.Prepend(vr)
				return vf, NewLazySequence(res), true
			}
		}
		return data.Nil, data.EmptyList, false
	}
	return NewLazySequence(res)
}
