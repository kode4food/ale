package special

import (
	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/generate"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Declare encodes global forward-declarations
func Declare(e encoder.Type, args ...data.Value) {
	arity.AssertMinimum(1, len(args))
	for _, v := range args {
		name := v.(data.LocalSymbol).Name()
		generate.Literal(e, name)
		e.Emit(isa.Declare)
	}
	if len(args) == 1 {
		generate.Literal(e, args[0])
	} else {
		v := data.NewVector(args...)
		generate.Literal(e, v)
	}
}

// Bind encodes a global definition
func Bind(e encoder.Type, args ...data.Value) {
	arity.AssertFixed(2, len(args))
	name := args[0].(data.LocalSymbol).Name()
	generate.Value(e, args[1])
	generate.Literal(e, name)
	e.Emit(isa.Bind)
	generate.Literal(e, args[0])
}
