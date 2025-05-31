package special

import (
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/eval"
	"github.com/kode4food/ale/pkg/macro"
)

type evalFunc func(env.Namespace, data.Value) (data.Value, error)

var (
	// Eval encodes an immediate evaluation
	Eval = makeEvaluator(eval.Value)

	// MacroExpand performs macro expansion of a form until it can no longer
	MacroExpand = makeEvaluator(macro.Expand)

	// MacroExpand1 performs a single-step macro expansion of a form
	MacroExpand1 = makeEvaluator(macro.Expand1)
)

func makeEvaluator(eval evalFunc) compiler.Call {
	return func(e encoder.Encoder, args ...data.Value) error {
		if err := data.CheckFixedArity(1, len(args)); err != nil {
			return err
		}
		if err := generate.Value(e, args[0]); err != nil {
			return err
		}
		ns := e.Globals()
		fn := data.MakeProcedure(func(args ...data.Value) data.Value {
			res, err := eval(ns, args[0])
			if err != nil {
				panic(err)
			}
			return res
		}, 1)
		if err := generate.Literal(e, fn); err != nil {
			return err
		}
		e.Emit(isa.Call1)
		return nil
	}
}
