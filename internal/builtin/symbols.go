package builtin

import "gitlab.com/kode4food/ale/api"

const anonName = api.Name("anon")

// Sym instantiates a new symbol (not interned)
func Sym(args ...api.Value) api.Value {
	if s, ok := args[0].(api.Symbol); ok {
		return s
	}
	s := args[0].(api.String)
	return api.ParseSymbol(api.Name(s))
}

// GenSym generates a unique symbol
func GenSym(args ...api.Value) api.Value {
	if len(args) == 0 {
		return api.NewGeneratedSymbol(anonName)
	}
	s := args[0].(api.String)
	return api.NewGeneratedSymbol(api.Name(s))
}

// IsSymbol returns whether or not the provided value is a symbol
func IsSymbol(args ...api.Value) api.Value {
	_, ok := args[0].(api.Symbol)
	return api.Bool(ok)
}

// IsLocal returns whether or not the provided value is an unqualified symbol
func IsLocal(args ...api.Value) api.Value {
	_, ok := args[0].(api.LocalSymbol)
	return api.Bool(ok)
}

// IsQualified returns whether or not the provided value is a qualified symbol
func IsQualified(args ...api.Value) api.Value {
	_, ok := args[0].(api.QualifiedSymbol)
	return api.Bool(ok)
}
