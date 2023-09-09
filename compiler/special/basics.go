package special

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/macro"
	"github.com/kode4food/ale/runtime/isa"
)

type evalFunc func(env.Namespace, data.Value) data.Value

// Eval encodes an immediate evaluation
func Eval(e encoder.Encoder, args ...data.Value) {
	data.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, makeEvaluator(e.Globals(), eval.Value))
	e.Emit(isa.Call1)
}

// MacroExpand performs macro expansion of a form until it can no longer
func MacroExpand(e encoder.Encoder, args ...data.Value) {
	data.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, makeEvaluator(e.Globals(), macro.Expand))
	e.Emit(isa.Call1)
}

// MacroExpand1 performs a single-step macro expansion of a form
func MacroExpand1(e encoder.Encoder, args ...data.Value) {
	data.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, makeEvaluator(e.Globals(), macro.Expand1))
	e.Emit(isa.Call1)
}

// Begin encodes a set of expressions, returning only the final evaluation
func Begin(e encoder.Encoder, args ...data.Value) {
	v := data.NewVector(args...)
	generate.Block(e, v)
}

func makeEvaluator(ns env.Namespace, eval evalFunc) data.Function {
	return data.Applicative(func(args ...data.Value) data.Value {
		return eval(ns, args[0])
	}, 1)
}
