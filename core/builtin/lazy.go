package builtin

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

// LazySequence treats a function as a lazy sequence
var LazySequence = data.MakeProcedure(func(args ...data.Value) data.Value {
	fn := args[0].(data.Procedure)
	resolver := makeLazyResolver(fn)
	return sequence.NewLazy(resolver)
}, 1)

func makeLazyResolver(f data.Procedure) sequence.LazyResolver {
	return func() (data.Value, data.Sequence, bool) {
		r := f.Call()
		if r != data.Null {
			s := r.(data.Sequence)
			if sf, sr, ok := s.Split(); ok {
				return sf, sr, true
			}
		}
		return data.Null, data.Null, false
	}
}
