package generate

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

var (
	vectorSym = env.RootSymbol("vector")
	objectSym = env.RootSymbol("object")
)

// Block encodes a set of expressions, returning only the final evaluation
func Block(e encoder.Encoder, s data.Sequence) error {
	f, r, ok := s.Split()
	if !ok {
		return Null(e)
	}
	if err := Value(e, f); err != nil {
		return err
	}
	for f, r, ok = r.Split(); ok; f, r, ok = r.Split() {
		e.Emit(isa.Pop)
		if err := Value(e, f); err != nil {
			return err
		}
	}
	return nil
}

// Vector encodes a vector
func Vector(e encoder.Encoder, v data.Vector) error {
	f, err := resolveBuiltIn(e, vectorSym)
	if err != nil {
		return err
	}
	return callStatic(e, f, v)
}

// Object encodes an object
func Object(e encoder.Encoder, a *data.Object) error {
	args := data.Vector{}
	for f, r, ok := a.Split(); ok; f, r, ok = r.Split() {
		v := f.(data.Pair)
		args = append(args, v.Car(), v.Cdr())
	}
	f, err := resolveBuiltIn(e, objectSym)
	if err != nil {
		return err
	}
	return callStatic(e, f, args)
}
