package builtin

import (
	"github.com/kode4food/ale/pkg/data"
)

var (
	// Bytes constructs a new byte vector
	Bytes = makeConstructor(data.NewBytes)

	// List constructs a new list
	List = makeConstructor(data.NewList)

	// Vector creates a new vector
	Vector = makeConstructor(data.NewVector)
)

func makeConstructor[T data.Value](orig func(...data.Value) T) data.Procedure {
	return data.MakeProcedure(func(args ...data.Value) data.Value {
		return orig(args...)
	})
}

func isPair(v data.Value) bool {
	p, ok := v.(data.Pair)
	return ok && p != data.Null
}
