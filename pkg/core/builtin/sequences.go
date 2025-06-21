package builtin

import (
	"errors"

	"github.com/kode4food/ale/pkg/data"
)

// ErrIndexOutOfBounds is raised when a call to Nth receives an index outside
// the bounds of the sequence being accessed
const ErrIndexOutOfBounds = "index out of bounds"

var (
	// Bytes constructs a new byte vector
	Bytes = makeConstructor(data.NewBytes)

	// List constructs a new list
	List = makeConstructor(data.NewList)

	// Vector creates a new vector
	Vector = makeConstructor(data.NewVector)
)

// Append adds a value to the end of the provided AppenderKey
var Append = data.MakeProcedure(func(args ...data.Value) data.Value {
	a := args[0].(data.Appender)
	s := args[1]
	return a.Append(s)
}, 2)

// Reverse returns a reversed copy of a Sequence
var Reverse = data.MakeProcedure(func(args ...data.Value) data.Value {
	r := args[0].(data.Reverser)
	return r.Reverse()
}, 1)

// Length returns the element count of the provided CountedKey
var Length = data.MakeProcedure(func(args ...data.Value) data.Value {
	s := args[0].(data.Counted)
	l := s.Count()
	return data.Integer(l)
}, 1)

// Nth returns the nth element of the provided sequence or a default
var Nth = data.MakeProcedure(func(args ...data.Value) data.Value {
	s := args[0].(data.Indexed)
	i := args[1].(data.Integer)
	if res, ok := s.ElementAt(int(i)); ok {
		return res
	}
	if len(args) > 2 {
		return args[2]
	}
	panic(errors.New(ErrIndexOutOfBounds))
}, 2, 3)

func makeConstructor[T data.Value](orig func(...data.Value) T) data.Procedure {
	return data.MakeProcedure(func(args ...data.Value) data.Value {
		return orig(args...)
	})
}

func isPair(v data.Value) bool {
	p, ok := v.(data.Pair)
	return ok && p != data.Null
}
