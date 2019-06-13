package generate

import (
	"fmt"

	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
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
func Block(e encoder.Type, s data.Sequence) {
	f, r, ok := s.Split()
	if !ok {
		Null(e)
		return
	}
	Value(e, f)
	for f, r, ok = r.Split(); ok; f, r, ok = r.Split() {
		e.Emit(isa.Pop)
		Value(e, f)
	}
}

// Sequence encodes a sequence
func Sequence(e encoder.Type, s data.Sequence) {
	switch typed := s.(type) {
	case data.NullType:
		Literal(e, typed)
	case data.String:
		Literal(e, typed)
	case data.List:
		Call(e, typed)
	case data.Vector:
		Vector(e, typed)
	case data.Associative:
		Associative(e, typed)
	default:
		panic(fmt.Errorf(CannotCompile, s))
	}
}

// List encodes a list
func List(e encoder.Type, l data.List) {
	f := resolveBuiltIn(e, listSym)
	args := stdlib.SequenceToValues(l)
	callApplicative(e, f.Caller(), args)
}

// Vector encodes a vector
func Vector(e encoder.Type, v data.Vector) {
	f := resolveBuiltIn(e, vectorSym)
	callApplicative(e, f.Caller(), data.Values(v))
}

// Associative encodes an associative array
func Associative(e encoder.Type, a data.Associative) {
	args := make(data.Values, a.Count()*2)
	var i int
	for f, r, ok := a.Split(); ok; f, r, ok = r.Split() {
		v := f.(data.Vector)
		args[i], _ = v.ElementAt(0)
		args[i+1], _ = v.ElementAt(1)
		i += 2
	}
	f := resolveBuiltIn(e, assocSym)
	callApplicative(e, f.Caller(), args)
}

func resolveBuiltIn(e encoder.Type, sym data.Symbol) data.Caller {
	manager := e.Globals().Manager()
	root := manager.GetRoot()
	res := namespace.MustResolveValue(root, sym)
	return res.(data.Caller)
}
