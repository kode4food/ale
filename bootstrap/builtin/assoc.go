package builtin

import "gitlab.com/kode4food/ale/data"

// Assoc creates a new associative array
func Assoc(args ...data.Value) data.Value {
	return data.NewAssociative(args...)
}

// IsAssoc returns whether or not a value is an associative array
func IsAssoc(args ...data.Value) data.Value {
	_, ok := args[0].(data.Associative)
	return data.Bool(ok)
}

// IsMapped returns whether or not a value is a mapped sequence
func IsMapped(args ...data.Value) data.Value {
	_, ok := args[0].(data.MappedSequence)
	return data.Bool(ok)
}
