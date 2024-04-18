package core

import "github.com/kode4food/ale/pkg/data"

const anonName = data.Local("anon")

// Sym instantiates a new symbol
var Sym = data.MakeProcedure(func(args ...data.Value) data.Value {
	if s, ok := args[0].(data.Symbol); ok {
		return s
	}
	s := args[0].(data.String)
	return data.MustParseSymbol(s)
}, 1)

// GenSym generates a unique symbol
var GenSym = data.MakeProcedure(func(args ...data.Value) data.Value {
	if len(args) == 0 {
		return data.NewGeneratedSymbol(anonName)
	}
	s := args[0].(data.Local)
	return data.NewGeneratedSymbol(s)
}, 0, 1)
