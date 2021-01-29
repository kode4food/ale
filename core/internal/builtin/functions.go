package builtin

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

// Apply performs a parameterized function call
var Apply = data.Applicative(func(args ...data.Value) data.Value {
	fn := args[0].(data.Function)
	al := len(args)
	if al == 2 {
		return fn.Call(sequence.ToValues(args[1].(data.Sequence))...)
	}
	last := al - 1
	ls := sequence.ToValues(args[last].(data.Sequence))
	prependedArgs := append(args[1:last], ls...)
	return fn.Call(prependedArgs...)
}, 2, data.OrMore)

// IsApply tests whether a value is callable
var IsApply = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Function)
	return data.Bool(ok)
}, 1)

// IsSpecial tests whether not a function is a special form
var IsSpecial = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(encoder.Call)
	return data.Bool(ok)
}, 1)
