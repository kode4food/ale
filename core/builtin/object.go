package builtin

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

// Object creates a new object instance
var Object = data.MakeProcedure(func(args ...ale.Value) ale.Value {
	res, err := data.ValuesToObject(args...)
	if err != nil {
		panic(err)
	}
	return res
})
