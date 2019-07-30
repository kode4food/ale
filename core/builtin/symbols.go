package builtin

import "github.com/kode4food/ale/data"

const anonName = data.Name("anon")

// Sym instantiates a new symbol
func Sym(args ...data.Value) data.Value {
	if s, ok := args[0].(data.Symbol); ok {
		return s
	}
	s := args[0].(data.String)
	return data.ParseSymbol(s)
}

// GenSym generates a unique symbol
func GenSym(args ...data.Value) data.Value {
	if len(args) == 0 {
		return data.NewGeneratedSymbol(anonName)
	}
	s := args[0].(data.String)
	return data.NewGeneratedSymbol(data.Name(s))
}

// IsSymbol returns whether or not the provided value is a symbol
func IsSymbol(args ...data.Value) data.Value {
	_, ok := args[0].(data.Symbol)
	return data.Bool(ok)
}

// IsLocal returns whether or not the provided value is an unqualified symbol
func IsLocal(args ...data.Value) data.Value {
	_, ok := args[0].(data.LocalSymbol)
	return data.Bool(ok)
}

// IsQualified returns whether or not the provided value is a qualified symbol
func IsQualified(args ...data.Value) data.Value {
	_, ok := args[0].(data.QualifiedSymbol)
	return data.Bool(ok)
}
