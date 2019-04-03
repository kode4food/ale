package builtin

import "gitlab.com/kode4food/ale/api"

// Vector creates a new vector
func Vector(args ...api.Value) api.Value {
	return api.Vector(args)
}

// IsVector returns whether or not the provided value is a vector
func IsVector(args ...api.Value) api.Value {
	_, ok := args[0].(api.Vector)
	return api.Bool(ok)
}
