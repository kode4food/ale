package generate

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/runtime/isa"
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
	l := len(v)
	switch {
	case l == 0:
		return Literal(e, data.EmptyVector)
	case l <= int(isa.OperandMask):
		ae := makeArgs(e, v)
		al, err := ae()
		if err != nil {
			return err
		}
		e.Emit(isa.Vector, isa.Operand(al))
		return nil
	default:
		f, err := resolveBuiltIn(e, vectorSym)
		if err != nil {
			return err
		}
		return callStatic(e, f, v)
	}
}

// Object encodes an object
func Object(e encoder.Encoder, a *data.Object) error {
	if a.IsEmpty() {
		return Literal(e, data.EmptyObject)
	}
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
