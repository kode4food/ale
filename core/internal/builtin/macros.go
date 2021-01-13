package builtin

import (
	"fmt"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/macro"
)

// Error messages
const (
	errCallableRequired = "argument must be callable: %s"
)

// Macro converts a function into a macro
func Macro(args ...data.Value) data.Value {
	switch arg0 := args[0].(type) {
	case data.Function:
		body := arg0.Call()
		wrapper := func(_ env.Namespace, args ...data.Value) data.Value {
			if err := arg0.CheckArity(len(args)); err != nil {
				panic(err)
			}
			return body(args...)
		}
		return macro.Call(wrapper)
	case data.Caller:
		body := arg0.Call()
		wrapper := func(_ env.Namespace, args ...data.Value) data.Value {
			return body(args...)
		}
		return macro.Call(wrapper)
	default:
		panic(fmt.Errorf(errCallableRequired, args[0]))
	}
}

// IsMacro returns whether the argument is a macro
func IsMacro(args ...data.Value) data.Value {
	_, ok := args[0].(macro.Call)
	return data.Bool(ok)
}
