package generate

import (
	"fmt"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/compiler/encoder"
	"gitlab.com/kode4food/ale/internal/namespace"
	"gitlab.com/kode4food/ale/internal/runtime/isa"
	"gitlab.com/kode4food/ale/stdlib"
)

// Error messages
const (
	CannotCompile = "sequence cannot be compiled: %s"
)

var (
	listSym   = namespace.RootSymbol("list")
	vectorSym = namespace.RootSymbol("vector")
	assocSym  = namespace.RootSymbol("assoc")
)

// Block encodes a set of expressions, returning only the final evaluation
func Block(e encoder.Type, s api.Sequence) {
	f, r, ok := s.Split()
	if !ok {
		Nil(e)
		return
	}
	Value(e, f)
	for f, r, ok = r.Split(); ok; f, r, ok = r.Split() {
		e.Append(isa.Pop)
		Value(e, f)
	}
}

// Sequence encodes a sequence
func Sequence(e encoder.Type, s api.Sequence) {
	switch typed := s.(type) {
	case api.String:
		Literal(e, typed)
	case *api.List:
		Call(e, typed)
	case api.Vector:
		Vector(e, typed)
	case api.Associative:
		Associative(e, typed)
	default:
		panic(fmt.Errorf(CannotCompile, s))
	}
}

// List encodes a list
func List(e encoder.Type, l *api.List) {
	f := resolveBuiltIn(e, listSym)
	args := stdlib.SequenceToValues(l)
	callApplicative(e, f.Call, args)
}

// Vector encodes a vector
func Vector(e encoder.Type, v api.Vector) {
	f := resolveBuiltIn(e, vectorSym)
	callApplicative(e, f.Call, api.Values(v))
}

// Associative encodes an associative array
func Associative(e encoder.Type, a api.Associative) {
	args := make(api.Values, a.Count()*2)
	var i int
	for f, r, ok := a.Split(); ok; f, r, ok = r.Split() {
		v := f.(api.Vector)
		args[i], _ = v.ElementAt(0)
		args[i+1], _ = v.ElementAt(1)
		i += 2
	}
	f := resolveBuiltIn(e, assocSym)
	callApplicative(e, f.Call, args)
}

func resolveBuiltIn(e encoder.Type, sym api.Symbol) *api.Function {
	manager := e.Globals().Manager()
	root := manager.GetRoot()
	res := namespace.MustResolveSymbol(root, sym)
	return res.(*api.Function)
}
