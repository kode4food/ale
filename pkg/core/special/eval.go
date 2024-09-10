package special

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/eval"
	"github.com/kode4food/ale/pkg/macro"
)

type evalFunc func(env.Namespace, data.Value) data.Value

var (
	// Eval encodes an immediate evaluation
	Eval = makeEvaluator(eval.Value)

	// MacroExpand performs macro expansion of a form until it can no longer
	MacroExpand = makeEvaluator(macro.Expand)

	// MacroExpand1 performs a single-step macro expansion of a form
	MacroExpand1 = makeEvaluator(macro.Expand1)
)

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
