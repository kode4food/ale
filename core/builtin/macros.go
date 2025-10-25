package builtin

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/macro"
)

// ErrProcedureRequired is raised when a call to the Macro built-in doesn't
// receive a data.Procedure to wrap
var ErrProcedureRequired = errors.New("argument must be a procedure")

// Macro converts a function into a macro
var Macro = data.MakeProcedure(func(args ...ale.Value) ale.Value {
	switch body := args[0].(type) {
	case data.Procedure:
		wrapper := func(_ env.Namespace, args ...ale.Value) ale.Value {
			if err := body.CheckArity(len(args)); err != nil {
				panic(err)
			}
			return body.Call(args...)
		}
		return macro.Call(wrapper)
	default:
		panic(fmt.Errorf("%w: %s", ErrProcedureRequired, args[0]))
	}
}, 1)

func isAtom(v ale.Value) bool {
	return !compiler.IsEvaluable(v)
}
