package core

import (
	"github.com/kode4food/ale/pkg/compiler/encoder"
	"github.com/kode4food/ale/pkg/compiler/generate"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/eval"
	"github.com/kode4food/ale/pkg/runtime/isa"
)

type evalFunc func(env.Namespace, data.Value) data.Value

// Eval encodes an immediate evaluation
var Eval = makeEvaluator(eval.Value)

func makeEvaluator(eval evalFunc) func(encoder.Encoder, ...data.Value) {
	return func(e encoder.Encoder, args ...data.Value) {
		data.AssertFixed(1, len(args))
		generate.Value(e, args[0])
		ns := e.Globals()
		fn := data.MakeProcedure(func(args ...data.Value) data.Value {
			return eval(ns, args[0])
		}, 1)
		generate.Literal(e, fn)
		e.Emit(isa.Call1)
	}
}
