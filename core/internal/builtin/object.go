package builtin

import (
	"errors"

	"github.com/kode4food/ale/data"
)

// Object creates a new object instance
var Object = data.Applicative(func(args ...data.Value) data.Value {
	res, err := data.ValuesToObject(args...)
	if err != nil {
		panic(err)
	}
	return res
})

// Get returns a value by key from the provided Mapper
var Get = data.Applicative(func(args ...data.Value) data.Value {
	s := args[0].(data.Mapped)
	res, _ := s.Get(args[1])
	return res
}, 2)

// Assoc returns a new Mapper containing the key/value association
var Assoc = data.Applicative(func(args ...data.Value) data.Value {
	s := args[0].(data.Mapper)
	if len(args) == 3 {
		p := data.NewCons(args[1], args[2])
		return s.Put(p)
	}
	if p, ok := args[1].(data.Pair); ok {
		return s.Put(p)
	}
	panic(errors.New(ErrPutRequiresPair))
}, 2, 3)

// Dissoc returns a new Mapper with the key removed
var Dissoc = data.Applicative(func(args ...data.Value) data.Value {
	s := args[0].(data.Mapper)
	if _, r, ok := s.Remove(args[1]); ok {
		return r
	}
	return s
}, 2)
