package builtin

import (
	"fmt"

	"github.com/kode4food/ale/compiler"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/macro"
)

// ErrProcedureRequired is raised when a call to the Macro built-in doesn't
// receive a data.Procedure to wrap
const ErrProcedureRequired = "argument must be a procedure: %s"

// MacroExpand performs macro expansion of a form until it can no longer
var MacroExpand = makeEvaluator(macro.Expand)

// MacroExpand1 performs a single-step macro expansion of a form
var MacroExpand1 = makeEvaluator(macro.Expand1)

// Macro converts a function into a macro
var Macro = data.MakeProcedure(func(args ...data.Value) data.Value {
	switch body := args[0].(type) {
	case data.Procedure:
		wrapper := func(_ env.Namespace, args ...data.Value) data.Value {
			if err := body.CheckArity(len(args)); err != nil {
				panic(err)
			}
			return body.Call(args...)
		}
		return macro.Call(wrapper)
	default:
		panic(fmt.Errorf(ErrProcedureRequired, args[0]))
	}
}, 1)

func isAtom(v data.Value) bool {
	return !compiler.IsEvaluable(v)
}
