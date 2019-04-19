package builtin

import "gitlab.com/kode4food/ale/data"

// Vector creates a new vector
func Vector(args ...data.Value) data.Value {
	return data.Vector(args)
}

// IsVector returns whether or not the provided value is a vector
func IsVector(args ...data.Value) data.Value {
	_, ok := args[0].(data.Vector)
	return data.Bool(ok)
}
