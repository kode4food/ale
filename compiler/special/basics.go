package special

import (
	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/generate"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Eval encodes an evaluation
func Eval(e encoder.Type, args ...data.Value) {
	arity.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, evalFor(e.Globals()))
	e.Emit(isa.Call1)
}

func evalFor(ns namespace.Type) data.Call {
	return data.Call(func(args ...data.Value) data.Value {
		return eval.Value(ns, args[0])
	})
}

// Begin encodes a set of expressions, returning only the final evaluation
func Begin(e encoder.Type, args ...data.Value) {
	v := data.NewVector(args...)
	generate.Block(e, v)
}
