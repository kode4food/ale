package builtin

import (
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/stdlib"
)

func makeLazyResolver(f data.Call) stdlib.LazyResolver {
	return func() (data.Value, data.Sequence, bool) {
		r := f()
		if r != data.Nil {
			s := r.(data.Sequence)
			if sf, sr, ok := s.Split(); ok {
				return sf, sr, true
			}
		}
		return data.Nil, data.EmptyList, false
	}
}

// LazySequence treats a function as a lazy sequence
func LazySequence(args ...data.Value) data.Value {
	fn := args[0].(data.Caller)
	resolver := makeLazyResolver(fn.Caller())
	return stdlib.NewLazySequence(resolver)
}

// Concat creates a lazy sequence that concatenates the provided sequences
func Concat(args ...data.Value) data.Value {
	switch len(args) {
	case 0:
		return data.EmptyList
	case 1:
		return args[0].(data.Sequence)
	default:
		return stdlib.Concat(args...)
	}
}

// Filter creates a lazy sequence that filters the provided sequence
func Filter(args ...data.Value) data.Value {
	fn := args[0].(data.Caller)
	s := args[1].(data.Sequence)
	return stdlib.Filter(s, fn.Caller())
}

// Map creates a lazy sequence that maps the provide sequence
func Map(args ...data.Value) data.Value {
	fn := args[0].(data.Caller)
	if len(args) == 2 {
		s := args[1].(data.Sequence)
		return stdlib.Map(s, fn.Caller())
	}
	return stdlib.MapParallel(data.Vector(args[1:]), fn.Caller())
}

// Reduce consumes the provided sequence, aggregating its values in some way
func Reduce(args ...data.Value) data.Value {
	fn := args[0].(data.Caller)
	if len(args) == 2 {
		s := args[1].(data.Sequence)
		return stdlib.Reduce(s, fn.Caller())
	}
	s := args[2].(data.Sequence).Prepend(args[1])
	return stdlib.Reduce(s, fn.Caller())
}

// ForEach calls a function for each element of the provided sequence
func ForEach(args ...data.Value) data.Value {
	seq := args[0].(data.Sequence)
	fn := args[1].(data.Caller).Caller()
	for f, r, ok := seq.Split(); ok; f, r, ok = r.Split() {
		fn(f)
	}
	return data.Nil
}
