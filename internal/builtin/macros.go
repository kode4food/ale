package builtin

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/macro"
)

// IsMacro returns whether or not the argument is a macro
func IsMacro(args ...api.Value) api.Value {
	if _, ok := args[0].(macro.Call); ok {
		return api.True
	}
	return api.False
}
