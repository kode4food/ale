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

// Get returns a value by key from the provided MappedSequence
var Get = data.Applicative(func(args ...data.Value) data.Value {
	s := args[0].(data.MappedSequence)
	res, _ := s.Get(args[1])
	return res
}, 2)

// Assoc returns a new MappedSequence containing the key/value association
var Assoc = data.Applicative(func(args ...data.Value) data.Value {
	s := args[0].(data.MappedSequence)
	if len(args) == 3 {
		p := data.NewCons(args[1], args[2])
		return s.Put(p)
	}
	if p, ok := args[1].(data.Pair); ok {
		return s.Put(p)
	}
	panic(errors.New(ErrPutRequiresPair))
}, 2, 3)

// Dissoc returns a new MappedSequence with the key removed
var Dissoc = data.Applicative(func(args ...data.Value) data.Value {
	s := args[0].(data.MappedSequence)
	if _, r, ok := s.Remove(args[1]); ok {
		return r
	}
	return s
}, 2)

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
