package builtin

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/stdlib"
)

func getCall(v api.Value) api.Call {
	return v.(api.Caller).Caller()
}

// Apply performs a parameterized function call
func Apply(args ...api.Value) api.Value {
	fn := getCall(args[0])
	if len(args) == 2 {
		return fn(stdlib.SequenceToVector(args[1].(api.Sequence))...)
	}
	last := len(args) - 1
	ls := stdlib.SequenceToVector(args[last].(api.Sequence))
	prependedArgs := append(args[1:last], ls...)
	return fn(prependedArgs...)
}

// Partial creates a partially bound version of a function
func Partial(args ...api.Value) api.Value {
	bound := getCall(args[0])
	values := args[1:]
	return bindFunction(bound, values)
}

func bindFunction(bound api.Call, outer api.Values) api.Call {
	return func(inner ...api.Value) api.Value {
		args := append(outer, inner...)
		return bound(args...)
	}
}

// IsApply tests whether or not a value is callable
func IsApply(args ...api.Value) api.Value {
	_, ok := args[0].(api.Caller)
	return api.Bool(ok)
}

// IsSpecial tests whether not a function is a special form
func IsSpecial(args ...api.Value) api.Value {
	if _, ok := args[0].(encoder.Call); ok {
		return api.True
	}
	return api.False
}
