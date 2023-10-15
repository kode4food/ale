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
var Eval = makeEvaluator(eval.Value)

// MacroExpand performs macro expansion of a form until it can no longer
var MacroExpand = makeEvaluator(macro.Expand)

// MacroExpand1 performs a single-step macro expansion of a form
var MacroExpand1 = makeEvaluator(macro.Expand1)

func makeEvaluator(eval evalFunc) func(encoder.Encoder, ...data.Value) {
	return func(e encoder.Encoder, args ...data.Value) {
		data.AssertFixed(1, len(args))
		generate.Value(e, args[0])
		ns := e.Globals()
		fn := data.Applicative(func(args ...data.Value) data.Value {
			return eval(ns, args[0])
		}, 1)
		generate.Literal(e, fn)
		e.Emit(isa.Call1)
	}
}

// Begin encodes a set of expressions, returning only the final evaluation
func Begin(e encoder.Encoder, args ...data.Value) {
	v := data.NewVector(args...)
	generate.Block(e, v)
}
