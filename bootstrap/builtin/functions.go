package builtin

import (
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/stdlib"
)

func getCall(v data.Value) data.Call {
	return v.(data.Caller).Caller()
}

// Apply performs a parameterized function call
func Apply(args ...data.Value) data.Value {
	fn := getCall(args[0])
	al := len(args)
	if al == 2 {
		return fn(stdlib.SequenceToVector(args[1].(data.Sequence))...)
	}
	last := al - 1
	ls := stdlib.SequenceToVector(args[last].(data.Sequence))
	prependedArgs := append(args[1:last], ls...)
	return fn(prependedArgs...)
}

// IsApply tests whether or not a value is callable
func IsApply(args ...data.Value) data.Value {
	if _, ok := args[0].(data.Caller); ok {
		return data.True
	}
	return data.False
}

// IsSpecial tests whether not a function is a special form
func IsSpecial(args ...data.Value) data.Value {
	if _, ok := args[0].(encoder.Call); ok {
		return data.True
	}
	return data.False
}
