package builtin

import (
	"errors"

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

// Assoc returns a new Mapper containing the key/value association
var Assoc = data.MakeProcedure(func(args ...data.Value) data.Value {
	s := args[0].(data.Mapper)
	for _, a := range args[1:] {
		if p, ok := a.(data.Pair); ok {
			s = s.Put(p).(data.Mapper)
			continue
		}
		panic(errors.New(ErrAssocRequiresPairs))
	}
	return s
}, 2, data.OrMore)

// Dissoc returns a new Mapper with the key removed
var Dissoc = data.MakeProcedure(func(args ...data.Value) data.Value {
	s := args[0].(data.Mapper)
	for _, k := range args[1:] {
		_, r, _ := s.Remove(k)
		s = r.(data.Mapper)
	}
	return s
}, 2, data.OrMore)
