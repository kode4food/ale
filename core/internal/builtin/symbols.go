package builtin

import "github.com/kode4food/ale/data"

const anonName = data.Local("anon")

// Sym instantiates a new symbol
var Sym = data.Applicative(func(args ...data.Value) data.Value {
	if s, ok := args[0].(data.Symbol); ok {
		return s
	}
	s := args[0].(data.String)
	return data.ParseSymbol(s)
}, 1)

// GenSym generates a unique symbol
var GenSym = data.Applicative(func(args ...data.Value) data.Value {
	if len(args) == 0 {
		return data.NewGeneratedSymbol(anonName)
	}
	s := args[0].(data.String)
	return data.NewGeneratedSymbol(data.Local(s))
}, 0, 1)

// IsSymbol returns whether the provided value is a symbol
var IsSymbol = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Symbol)
	return data.Bool(ok)
}, 1)

// IsLocal returns whether the provided value is an unqualified symbol
var IsLocal = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Local)
	return data.Bool(ok)
}, 1)

// IsQualified returns whether the provided value is a qualified symbol
var IsQualified = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Qualified)
	return data.Bool(ok)
}, 1)
