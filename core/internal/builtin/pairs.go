package builtin

import "github.com/kode4food/ale/data"

// Cons adds a value to the beginning of the provided Sequence
var Cons = data.Applicative(func(args ...data.Value) data.Value {
	car := args[0]
	cdr := args[1]
	if p, ok := cdr.(data.Prepender); ok {
		return p.Prepend(car)
	}
	return data.NewCons(car, cdr)
}, 2)
