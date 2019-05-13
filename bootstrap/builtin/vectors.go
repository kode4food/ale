package builtin

import "gitlab.com/kode4food/ale/data"

// Vector creates a new vector
func Vector(args ...data.Value) data.Value {
	return data.NewVector(args...)
}

// IsVector returns whether or not the provided value is a vector
func IsVector(args ...data.Value) data.Value {
	_, ok := args[0].(data.Vector)
	return data.Bool(ok)
}

// IsAppender returns whether or not the provided value is an appender
func IsAppender(args ...data.Value) data.Value {
	_, ok := args[0].(data.Appender)
	return data.Bool(ok)
}
