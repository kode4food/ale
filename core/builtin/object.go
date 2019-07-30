package builtin

import "github.com/kode4food/ale/data"

// Object creates a new object instance
func Object(args ...data.Value) data.Value {
	return data.NewObject(args...)
}

// IsObject returns whether or not a value is an object
func IsObject(args ...data.Value) data.Value {
	_, ok := args[0].(data.Object)
	return data.Bool(ok)
}

// IsMapped returns whether or not a value is a mapped sequence
func IsMapped(args ...data.Value) data.Value {
	_, ok := args[0].(data.MappedSequence)
	return data.Bool(ok)
}
