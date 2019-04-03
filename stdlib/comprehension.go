package stdlib

import (
	"sync"
	"sync/atomic"

	"gitlab.com/kode4food/ale/api"
)

// Map creates a new mapped Sequence
func Map(s api.Sequence, mapper api.Call) api.Sequence {
	var res LazyResolver
	next := s

	res = func() (api.Value, api.Sequence, bool) {
		if f, r, ok := next.Split(); ok {
			m := mapper(f)
			next = r
			return m, NewLazySequence(res), true
		}
		return api.Nil, api.EmptyList, false
	}
	return NewLazySequence(res)
}

// MapParallel creates a new mapped Sequence from a Sequence of Sequences
// that are used to provide multiple arguments to the mapper function
func MapParallel(s api.Sequence, mapper api.Call) api.Sequence {
	var res LazyResolver
	next := make([]api.Sequence, 0)
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		next = append(next, f.(api.Sequence))
	}
	nextLen := len(next)

	res = func() (api.Value, api.Sequence, bool) {
		var exhausted int32
		args := make(api.Vector, nextLen)

		var wg sync.WaitGroup
		wg.Add(nextLen)

		for i, s := range next {
			go func(i int, s api.Sequence) {
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
			return api.Nil, api.EmptyList, false
		}
		m := mapper(args...)
		return m, NewLazySequence(res), true
	}
	return NewLazySequence(res)
}

// Filter creates a new filtered Sequence
func Filter(s api.Sequence, filter api.Call) api.Sequence {
	var res LazyResolver
	next := s

	res = func() (api.Value, api.Sequence, bool) {
		for f, r, ok := next.Split(); ok; f, r, ok = r.Split() {
			next = r
			if api.Truthy(filter(f)) {
				return f, NewLazySequence(res), true
			}
		}
		return api.Nil, api.EmptyList, false
	}
	return NewLazySequence(res)
}

// Concat creates a new sequence based on the content of several Sequences
func Concat(s ...api.Value) api.Sequence {
	var res LazyResolver
	var next api.Sequence = api.Vector(s)

	res = func() (api.Value, api.Sequence, bool) {
		for f, r, ok := next.Split(); ok; f, r, ok = r.Split() {
			v := f.(api.Sequence)
			next = r
			if vf, vr, ok := v.Split(); ok {
				next = next.Prepend(vr)
				return vf, NewLazySequence(res), true
			}
		}
		return api.Nil, api.EmptyList, false
	}
	return NewLazySequence(res)
}

// Take creates a Sequence based on the first elements of the source
func Take(s api.Sequence, count api.Integer) api.Sequence {
	var res LazyResolver
	var idx api.Integer
	next := s

	res = func() (api.Value, api.Sequence, bool) {
		if f, r, ok := next.Split(); ok && idx < count {
			next = r
			idx++
			return f, NewLazySequence(res), true
		}
		return api.Nil, api.EmptyList, false
	}
	return NewLazySequence(res)
}

// Drop creates a Sequence based on dropping the first elements of the source
func Drop(s api.Sequence, count api.Integer) api.Sequence {
	var first, rest LazyResolver
	next := s

	first = func() (api.Value, api.Sequence, bool) {
		for i := api.Integer(0); i < count && next.IsSequence(); i++ {
			next = next.Rest()
		}
		return rest()
	}

	rest = func() (api.Value, api.Sequence, bool) {
		if f, r, ok := next.Split(); ok {
			next = r
			return f, NewLazySequence(rest), true
		}
		return api.Nil, api.EmptyList, false
	}

	return NewLazySequence(first)
}

// Reduce performs a reduce operation over a Sequence, starting with the
// first two elements of that sequence.
func Reduce(s api.Sequence, reduce api.Call) api.Value {
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
