package builtin

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/stdlib"
)

func makeLazyResolver(f api.Call) stdlib.LazyResolver {
	return func() (api.Value, api.Sequence, bool) {
		r := f()
		if r != api.Nil {
			s := r.(api.Sequence)
			if sf, sr, ok := s.Split(); ok {
				return sf, sr, true
			}
		}
		return api.Nil, api.EmptyList, false
	}
}

// LazySequence treats a function as a lazy sequence
func LazySequence(args ...api.Value) api.Value {
	fn := args[0].(api.Caller)
	resolver := makeLazyResolver(fn.Caller())
	return stdlib.NewLazySequence(resolver)
}

// Concat creates a lazy sequence that concatenates the provided sequences
func Concat(args ...api.Value) api.Value {
	switch len(args) {
	case 0:
		return api.EmptyList
	case 1:
		return args[0].(api.Sequence)
	default:
		return stdlib.Concat(args...)
	}
}

// Filter creates a lazy sequence that filters the provided sequence
func Filter(args ...api.Value) api.Value {
	fn := args[0].(api.Caller)
	s := args[1].(api.Sequence)
	return stdlib.Filter(s, fn.Caller())
}

// Map creates a lazy sequence that maps the provide sequence
func Map(args ...api.Value) api.Value {
	fn := args[0].(api.Caller)
	if len(args) == 2 {
		s := args[1].(api.Sequence)
		return stdlib.Map(s, fn.Caller())
	}
	return stdlib.MapParallel(api.Vector(args[1:]), fn.Caller())
}

// Take returns the first 'n' elements of the provided sequence
func Take(args ...api.Value) api.Value {
	n := args[0].(api.Integer)
	s := args[1].(api.Sequence)
	return stdlib.Take(s, n)
}

// Drop discards the first 'n' elements of the provided sequence
func Drop(args ...api.Value) api.Value {
	n := args[0].(api.Integer)
	s := args[1].(api.Sequence)
	return stdlib.Drop(s, n)
}

// Reduce consumes the provided sequence, aggregating its values in some way
func Reduce(args ...api.Value) api.Value {
	fn := args[0].(api.Caller)
	if len(args) == 2 {
		s := args[1].(api.Sequence)
		return stdlib.Reduce(s, fn.Caller())
	}
	s := args[2].(api.Sequence).Prepend(args[1])
	return stdlib.Reduce(s, fn.Caller())
}

// ForEach calls a function for each element of the provided sequence
func ForEach(args ...api.Value) api.Value {
	seq := args[0].(api.Sequence)
	fn := args[1].(api.Caller).Caller()
	for f, r, ok := seq.Split(); ok; f, r, ok = r.Split() {
		fn(f)
	}
	return api.Nil
}
