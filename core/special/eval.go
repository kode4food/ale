package special

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/macro"
)

type evalFunc func(env.Namespace, ale.Value) (ale.Value, error)

var (
	// Eval encodes an immediate evaluation
	Eval = makeEvaluator(eval.Value)

	// MacroExpand performs macro expansion of a form until it can no longer
	MacroExpand = makeEvaluator(macro.Expand)

	// MacroExpand1 performs a single-step macro expansion of a form
	MacroExpand1 = makeEvaluator(macro.Expand1)
)

func makeEvaluator(eval evalFunc) compiler.Call {
	return func(e encoder.Encoder, args ...ale.Value) error {
		if err := data.CheckFixedArity(1, len(args)); err != nil {
			return err
		}
		if err := generate.Value(e, args[0]); err != nil {
			return err
		}
		ns := e.Globals()
		fn := data.MakeProcedure(func(args ...ale.Value) ale.Value {
			res, err := eval(ns, args[0])
			if err != nil {
				panic(err)
			}
			return res
		})
		if err := generate.Literal(e, fn); err != nil {
			return err
		}
		e.Emit(isa.Call1)
		return nil
	}
}
