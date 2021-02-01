package builtin

import (
	"github.com/kode4food/ale/compiler"
	"github.com/kode4food/ale/data"
)

// IsIdentical returns whether the values represent the same object
var IsIdentical = data.Applicative(func(args ...data.Value) data.Value {
	l := args[0]
	for _, f := range args[1:] {
		if !l.Equal(f) {
			return data.False
		}
	}
	return data.True
}, 1, data.OrMore)

// IsAtom returns whether the provided value is atomic
var IsAtom = data.Applicative(func(args ...data.Value) data.Value {
	return data.Bool(!compiler.IsEvaluable(args[0]))
}, 1)

// IsBoolean returns whether the provided value is a boolean
var IsBoolean = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Bool)
	return data.Bool(ok)
}, 1)

// IsKeyword returns whether the provided value is a keyword
var IsKeyword = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Keyword)
	return data.Bool(ok)
}, 1)
