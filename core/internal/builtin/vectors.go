package builtin

import "github.com/kode4food/ale/data"

// Vector creates a new vector
var Vector = data.Applicative(func(args ...data.Value) data.Value {
	return data.NewVector(args...)
})

// IsVector returns whether the provided value is a vector
var IsVector = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Vector)
	return data.Bool(ok)
}, 1)

// IsAppender returns whether the provided value is an appender
var IsAppender = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.AppenderSequence)
	return data.Bool(ok)
}, 1)
