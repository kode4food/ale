package builtin

import (
	"fmt"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/macro"
	"github.com/kode4food/ale/namespace"
)

// Error messages
const (
	CallableRequired = "argument must be callable: %s"
)

// Macro converts a function into a macro
func Macro(args ...data.Value) data.Value {
	switch typed := args[0].(type) {
	case data.Function:
		body := typed.Call()
		wrapper := func(_ namespace.Type, args ...data.Value) data.Value {
			if err := typed.CheckArity(len(args)); err != nil {
				panic(err)
			}
			return body(args...)
		}
		return macro.Call(wrapper)
	case data.Caller:
		body := typed.Call()
		wrapper := func(_ namespace.Type, args ...data.Value) data.Value {
			return body(args...)
		}
		return macro.Call(wrapper)
	default:
		panic(fmt.Errorf(CallableRequired, args[0]))
	}
}

// IsMacro returns whether the argument is a macro
func IsMacro(args ...data.Value) data.Value {
	_, ok := args[0].(macro.Call)
	return data.Bool(ok)
}
