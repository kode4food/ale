package builtin

import (
	"fmt"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/macro"
)

// Error messages
const (
	ErrFunctionRequired = "argument must be a function: %s"
)

// Macro converts a function into a macro
var Macro = data.Applicative(func(args ...data.Value) data.Value {
	switch body := args[0].(type) {
	case data.Function:
		wrapper := func(_ env.Namespace, args ...data.Value) data.Value {
			if err := body.CheckArity(len(args)); err != nil {
				panic(err)
			}
			return body.Call(args...)
		}
		return macro.Call(wrapper)
	default:
		panic(fmt.Errorf(ErrFunctionRequired, args[0]))
	}
}, 1)

// IsMacro returns whether the argument is a macro
var IsMacro = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(macro.Call)
	return data.Bool(ok)
}, 1)
