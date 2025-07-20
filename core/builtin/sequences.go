package builtin

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

var (
	// Bytes constructs a new byte vector
	Bytes = makeConstructor(data.NewBytes)

	// List constructs a new list
	List = makeConstructor(data.NewList)

	// Vector creates a new vector
	Vector = makeConstructor(data.NewVector)
)

func makeConstructor[T ale.Value](orig func(...ale.Value) T) data.Procedure {
	return data.MakeProcedure(func(args ...ale.Value) ale.Value {
		return orig(args...)
	})
}

func isPair(v ale.Value) bool {
	p, ok := v.(data.Pair)
	return ok && p != data.Null
}
