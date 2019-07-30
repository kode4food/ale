package special

import (
	"github.com/kode4food/ale/compiler/arity"
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

// If encodes an (if predicate consequent alternative) form
func If(e encoder.Type, args ...data.Value) {
	al := arity.AssertRanged(2, 3, len(args))
	generate.Branch(e,
		func() {
			generate.Value(e, args[0])
			e.Emit(isa.MakeTruthy)
		},
		func() {
			generate.Value(e, args[1])
		},
		func() {
			if al == 3 {
				generate.Value(e, args[2])
			} else {
				generate.Null(e)
			}
		},
	)
}
