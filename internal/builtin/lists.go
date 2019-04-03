package builtin

import "gitlab.com/kode4food/ale/api"

// List constructs a new list
func List(args ...api.Value) api.Value {
	return api.NewList(args...)
}

// IsList returns whether or not the provided value is a list
func IsList(args ...api.Value) api.Value {
	_, ok := args[0].(*api.List)
	return api.Bool(ok)
}
