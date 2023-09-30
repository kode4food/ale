package builtin

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

// Call performs a function call
var Call = data.Applicative(func(args ...data.Value) data.Value {
	fn := args[0].(data.Function)
	v := sequence.ToValues(args[1].(data.Sequence))
	return fn.Call(v...)
}, 2)
