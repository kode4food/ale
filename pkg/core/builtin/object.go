package builtin

import (
	"github.com/kode4food/ale/pkg/data"
)

// ErrAssocRequiresPairs is raised when a call to Assoc receives an argument
// other than a Pair
const ErrAssocRequiresPairs = "assoc requires one or more pairs"

// Object creates a new object instance
var Object = data.MakeProcedure(func(args ...data.Value) data.Value {
	res, err := data.ValuesToObject(args...)
	if err != nil {
		panic(err)
	}
	return res
})

// Get returns a value by key from the provided Mapper
var Get = data.MakeProcedure(func(args ...data.Value) data.Value {
	s := args[0].(data.Mapped)
	res, _ := s.Get(args[1])
	return res
}, 2)
