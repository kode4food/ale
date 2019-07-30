package generate

import (
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/namespace"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/ale/stdlib"
)

// Error messages
const (
	CannotCompile = "sequence cannot be compiled: %s"
)

var (
	listSym   = namespace.RootSymbol("list")
	vectorSym = namespace.RootSymbol("vector")
	objectSym = namespace.RootSymbol("object")
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
	case data.String:
		Literal(e, typed)
	case data.List:
		Call(e, typed)
	case data.Vector:
		Vector(e, typed)
	case data.Object:
		Object(e, typed)
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

// Object encodes an object
func Object(e encoder.Type, a data.Object) {
	args := data.Values{}
	for f, r, ok := a.Split(); ok; f, r, ok = r.Split() {
		v := f.(data.Pair)
		args = append(args, v.Car(), v.Cdr())
	}
	f := resolveBuiltIn(e, objectSym)
	callApplicative(e, f.Caller(), args)
}
