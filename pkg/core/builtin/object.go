package builtin

import "github.com/kode4food/ale/pkg/data"

// Object creates a new object instance
var Object = data.MakeProcedure(func(args ...data.Value) data.Value {
	res, err := data.ValuesToObject(args...)
	if err != nil {
		panic(err)
	}
	return res
})
