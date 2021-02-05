package generate

import (
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/runtime/isa"
)

// Error messages
const (
	errCannotCompile = "sequence cannot be compiled: %s"
)

var (
	vectorSym = env.RootSymbol("vector")
	objectSym = env.RootSymbol("object")
)

// Block encodes a set of expressions, returning only the final evaluation
func Block(e encoder.Encoder, s data.Sequence) {
	f, r, ok := s.Split()
	if !ok {
		Nil(e)
		return
	}
	Value(e, f)
	for f, r, ok = r.Split(); ok; f, r, ok = r.Split() {
		e.Emit(isa.Pop)
		Value(e, f)
	}
}

// Sequence encodes a sequence
func Sequence(e encoder.Encoder, s data.Sequence) {
	switch s := s.(type) {
	case data.String:
		Literal(e, s)
	case data.List:
		Call(e, s)
	case data.Vector:
		Vector(e, s)
	case data.Object:
		Object(e, s)
	default:
		// Programmer error
		panic(fmt.Errorf(errCannotCompile, s))
	}
}

// Vector encodes a vector
func Vector(e encoder.Encoder, v data.Vector) {
	f := resolveBuiltIn(e, vectorSym)
	callFunction(e, f, data.Values(v))
}

// Object encodes an object
func Object(e encoder.Encoder, a data.Object) {
	args := data.Values{}
	for f, r, ok := a.Split(); ok; f, r, ok = r.Split() {
		v := f.(data.Pair)
		args = append(args, v.Car(), v.Cdr())
	}
	f := resolveBuiltIn(e, objectSym)
	callApplicative(e, f, args)
}
