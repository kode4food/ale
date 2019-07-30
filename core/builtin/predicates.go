package builtin

import (
	"github.com/kode4food/ale/compiler"
	"github.com/kode4food/ale/data"
)

// IsIdentical returns whether or not the two values represent the same object
func IsIdentical(args ...data.Value) data.Value {
	l := args[0]
	for _, f := range args[1:] {
		if l != f {
			return data.False
		}
	}
	return data.True
}

// IsAtom returns whether or not the provided value is atomic
func IsAtom(args ...data.Value) data.Value {
	return data.Bool(!compiler.IsEvaluable(args[0]))
}

// IsBoolean returns whether or not the provided value is a boolean
func IsBoolean(args ...data.Value) data.Value {
	_, ok := args[0].(data.Bool)
	return data.Bool(ok)
}

// IsKeyword returns whether or not the provided value is a keyword
func IsKeyword(args ...data.Value) data.Value {
	_, ok := args[0].(data.Keyword)
	return data.Bool(ok)
}
