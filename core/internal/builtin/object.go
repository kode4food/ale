package builtin

import "github.com/kode4food/ale/data"

// Object creates a new object instance
var Object = data.Applicative(func(args ...data.Value) data.Value {
	res, err := data.ValuesToObject(args...)
	if err != nil {
		panic(err)
	}
	return res
})

// IsObject returns whether a value is an object
var IsObject = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Object)
	return data.Bool(ok)
}, 1)

// IsMapped returns whether a value is a mapped sequence
var IsMapped = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.MappedSequence)
	return data.Bool(ok)
}, 1)
