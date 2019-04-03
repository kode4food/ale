package builtin

import "gitlab.com/kode4food/ale/api"

// IsMacro returns whether or not the argument is a macro
func IsMacro(args ...api.Value) api.Value {
	if f, ok := args[0].(*api.Function); ok {
		return api.Bool(f.IsMacro())
	}
	return api.False
}
