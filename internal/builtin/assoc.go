package builtin

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/stdlib"
)

// Assoc creates a new associative array
func Assoc(args ...api.Value) api.Value {
	return stdlib.SequenceToAssociative(api.Vector(args))
}

// IsAssoc returns whether or not a value is an associative array
func IsAssoc(args ...api.Value) api.Value {
	_, ok := args[0].(api.Associative)
	return api.Bool(ok)
}

// IsMapped returns whether or not a value is a mapped sequence
func IsMapped(args ...api.Value) api.Value {
	_, ok := args[0].(api.MappedSequence)
	return api.Bool(ok)
}
