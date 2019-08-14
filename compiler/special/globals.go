package special

import (
	"github.com/kode4food/ale/compiler/arity"
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

// Declare encodes a global forward-declaration
func Declare(e encoder.Type, args ...data.Value) {
	arity.AssertFixed(1, len(args))
	name := args[0].(data.LocalSymbol).Name()
	generate.Literal(e, name)
	e.Emit(isa.Declare)
	generate.Literal(e, args[0])
}

// Define encodes a global definition
func Define(e encoder.Type, args ...data.Value) {
	arity.AssertFixed(2, len(args))
	name := args[0].(data.LocalSymbol).Name()
	generate.Value(e, args[1])
	generate.Literal(e, name)
	e.Emit(isa.Bind)
	generate.Literal(e, args[0])
}
