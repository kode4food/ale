package builtin

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/compiler"
)

// IsIdentical returns whether or not the two values represent the same object
func IsIdentical(args ...api.Value) api.Value {
	l := args[0]
	for _, f := range args[1:] {
		if l != f {
			return api.False
		}
	}
	return api.True
}

// IsAtom returns whether or not the provided value is atomic
func IsAtom(args ...api.Value) api.Value {
	if compiler.IsEvaluable(args[0]) {
		return api.False
	}
	return api.True
}

// IsNil returns whether or not the provided value is nil
func IsNil(args ...api.Value) api.Value {
	return api.Bool(args[0] == api.Nil)
}

// IsKeyword returns whether or not the provided value is a keyword
func IsKeyword(args ...api.Value) api.Value {
	_, ok := args[0].(api.Keyword)
	return api.Bool(ok)
}
