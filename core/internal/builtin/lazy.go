package builtin

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

func makeLazyResolver(f data.Call) sequence.LazyResolver {
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
	resolver := makeLazyResolver(fn.Call())
	return sequence.NewLazy(resolver)
}
