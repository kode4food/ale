package core

import (
	"fmt"

	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/macro"
)

// ErrProcedureRequired is raised when a call to the Macro built-in doesn't
// receive a data.Procedure to wrap
const ErrProcedureRequired = "argument must be a procedure: %s"

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
