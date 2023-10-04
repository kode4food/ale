package builtin

import "github.com/kode4food/ale/data"

// IsNaN returns true if the provided value is not a number
var IsNaN = data.Applicative(func(args ...data.Value) data.Value {
	if num, ok := args[0].(data.Number); ok {
		return data.Bool(num.IsNaN())
	}
	return data.False
}, 1)
