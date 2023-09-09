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
