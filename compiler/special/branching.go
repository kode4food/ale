package special

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

// If encodes an (if predicate consequent alternative) form
func If(e encoder.Encoder, args ...data.Value) {
	al := data.AssertRanged(2, 3, len(args))
	branch(e,
		data.Applicative(func(...data.Value) data.Value {
			value(e, args[0])
			e.Emit(isa.MakeTruthy)
			return data.Nil
		}, 0),
		data.Applicative(func(...data.Value) data.Value {
			value(e, args[1])
			return data.Nil
		}, 0),
		data.Applicative(func(...data.Value) data.Value {
			if al == 3 {
				value(e, args[2])
			} else {
				e.Emit(isa.Nil)
			}
			return data.Nil
		}, 0),
	)
}
