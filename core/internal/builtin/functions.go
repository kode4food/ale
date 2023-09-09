package builtin

import (
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
