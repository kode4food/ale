package special

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

// Declare encodes a public global forward-declaration
func Declare(e encoder.Encoder, args ...data.Value) {
	declare(e, isa.Declare, args...)
}

// Private encodes a private global forward-declaration
func Private(e encoder.Encoder, args ...data.Value) {
	declare(e, isa.Private, args...)
}

func declare(e encoder.Encoder, oc isa.Opcode, args ...data.Value) {
	data.AssertFixed(1, len(args))
	generate.Literal(e, args[0].(data.LocalSymbol))
	e.Emit(isa.Dup)
	e.Emit(oc)
}

// Define encodes a global definition
func Define(e encoder.Encoder, args ...data.Value) {
	data.AssertFixed(2, len(args))
	generate.Literal(e, args[0].(data.LocalSymbol))
	e.Emit(isa.Dup)
	generate.Value(e, args[1])
	e.Emit(isa.Bind)
}
