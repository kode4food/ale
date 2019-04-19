package builtin

import (
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/macro"
)

// IsMacro returns whether or not the argument is a macro
func IsMacro(args ...data.Value) data.Value {
	if _, ok := args[0].(macro.Call); ok {
		return data.True
	}
	return data.False
}
