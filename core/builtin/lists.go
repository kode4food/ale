package builtin

import "gitlab.com/kode4food/ale/data"

// List constructs a new list
func List(args ...data.Value) data.Value {
	return data.NewList(args...)
}

// IsList returns whether or not the provided value is a list
func IsList(args ...data.Value) data.Value {
	_, ok := args[0].(*data.List)
	return data.Bool(ok)
}
