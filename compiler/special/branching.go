package special

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

// If encodes an (if predicate consequent alternative) form
func If(e encoder.Encoder, args ...data.Value) {
	al := data.AssertRanged(2, 3, len(args))
	generate.Branch(e,
		func(e encoder.Encoder) {
			generate.Value(e, args[0])
			e.Emit(isa.MakeTruthy)
		},
		func(e encoder.Encoder) {
			generate.Value(e, args[1])
		},
		func(e encoder.Encoder) {
			if al == 3 {
				generate.Value(e, args[2])
			} else {
				e.Emit(isa.Nil)
			}
		},
	)
}
