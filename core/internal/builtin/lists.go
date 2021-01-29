package builtin

import "github.com/kode4food/ale/data"

// List constructs a new list
var List = data.Applicative(func(args ...data.Value) data.Value {
	return data.NewList(args...)
})

// IsList returns whether the provided value is a list
var IsList = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.List)
	return data.Bool(ok)
}, 1)
