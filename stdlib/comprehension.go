package stdlib

import (
	"sync"
	"sync/atomic"

	"gitlab.com/kode4food/ale/data"
)

// Map creates a new mapped Sequence
func Map(s data.Sequence, mapper data.Call) data.Sequence {
	var res LazyResolver
	next := s

	res = func() (data.Value, data.Sequence, bool) {
		if f, r, ok := next.Split(); ok {
			m := mapper(f)
			next = r
			return m, NewLazySequence(res), true
		}
		return data.Nil, data.EmptyList, false
	}
	return NewLazySequence(res)
}

// MapParallel creates a new mapped Sequence from a Sequence of Sequences
// that are used to provide multiple arguments to the mapper function
func MapParallel(s data.Sequence, mapper data.Call) data.Sequence {
	var res LazyResolver
	next := make([]data.Sequence, 0)
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		next = append(next, f.(data.Sequence))
	}
	nextLen := len(next)

	res = func() (data.Value, data.Sequence, bool) {
		var exhausted int32
		args := make(data.Vector, nextLen)

		var wg sync.WaitGroup
		wg.Add(nextLen)

		for i, s := range next {
			go func(i int, s data.Sequence) {
				if f, r, ok := s.Split(); ok {
					args[i] = f
					next[i] = r
				} else {
					atomic.StoreInt32(&exhausted, 1)
				}
				wg.Done()
			}(i, s)
		}
		wg.Wait()

		if exhausted > 0 {
			return data.Nil, data.EmptyList, false
		}
		m := mapper(args...)
		return m, NewLazySequence(res), true
	}
	return NewLazySequence(res)
}

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
	var next data.Sequence = data.Vector(s)

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

// Reduce performs a reduce operation over a Sequence, starting with the
// first two elements of that sequence.
func Reduce(s data.Sequence, reduce data.Call) data.Value {
	arg1, r, ok := s.Split()
	if !ok {
		return reduce()
	}
	arg2, r, ok := r.Split()
	if !ok {
		return reduce(arg1)
	}
	res := reduce(arg1, arg2)
	for f, r, ok := r.Split(); ok; f, r, ok = r.Split() {
		res = reduce(res, f)
	}
	return res
}
