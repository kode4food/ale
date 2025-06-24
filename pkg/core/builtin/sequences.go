package builtin

import (
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

func makeConstructor[T data.Value](orig func(...data.Value) T) data.Procedure {
	return data.MakeProcedure(func(args ...data.Value) data.Value {
		return orig(args...)
	})
}

func isPair(v data.Value) bool {
	p, ok := v.(data.Pair)
	return ok && p != data.Null
}
