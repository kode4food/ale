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