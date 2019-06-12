package builtin

import (
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/macro"
)

// IsMacro returns whether or not the argument is a macro
func IsMacro(args ...data.Value) data.Value {
	_, ok := args[0].(macro.Call)
	return data.Bool(ok)
}
