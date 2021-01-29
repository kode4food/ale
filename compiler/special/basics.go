package special

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/runtime/isa"
)

// Eval encodes an immediate evaluation
func Eval(e encoder.Encoder, args ...data.Value) {
	data.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, evalFor(e.Globals()))
	e.Emit(isa.Call1)
}

func evalFor(ns env.Namespace) data.Function {
	return data.Applicative(func(args ...data.Value) data.Value {
		return eval.Value(ns, args[0])
	}, 1)
}

// Begin encodes a set of expressions, returning only the final evaluation
func Begin(e encoder.Encoder, args ...data.Value) {
	v := data.NewVector(args...)
	generate.Block(e, v)
}
