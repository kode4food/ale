package special

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/runtime/isa"
)

var (
	ErrExpectedName     = errors.New("name expected")
	ErrUnexpectedImport = errors.New("unexpected import pattern")
	ErrDuplicateName    = errors.New("duplicate name in import")
)

func MakeNamespace(e encoder.Encoder, args ...ale.Value) error {
	if err := data.CheckFixedArity(1, len(args)); err != nil {
		return err
	}
	name, ok := args[0].(data.Local)
	if !ok {
		return fmt.Errorf("%w: %s", ErrExpectedName, args[0])
	}
	fn := data.MakeProcedure(func(...ale.Value) ale.Value {
		ns, err := e.Globals().Environment().NewQualified(name)
		if err != nil {
			panic(err)
		}
		return makeInNamespaceCall(ns)
	})
	if err := generate.Literal(e, fn); err != nil {
		return err
	}
	e.Emit(isa.Call0)
	return nil
}

func makeInNamespaceCall(ns env.Namespace) compiler.Call {
	return func(e encoder.Encoder, args ...ale.Value) error {
		fn := data.MakeProcedure(func(...ale.Value) ale.Value {
			res, err := eval.Block(ns, data.Vector(args))
			if err != nil {
				panic(err)
			}
			return res
		})
		if err := generate.Literal(e, fn); err != nil {
			return err
		}
		e.Emit(isa.Call0)
		return nil
	}
}

func Declared(e encoder.Encoder, args ...ale.Value) error {
	if err := data.CheckRangedArity(0, 1, len(args)); err != nil {
		return err
	}
	ns := e.Globals()
	if len(args) > 0 {
		name, ok := args[0].(data.Local)
		if !ok {
			return fmt.Errorf("%w: %s", ErrExpectedName, args[0])
		}
		var err error
		ns, err = ns.Environment().GetQualified(name)
		if err != nil {
			return err
		}
	}
	fn := data.MakeProcedure(func(...ale.Value) ale.Value {
		return localsToVector(ns.Declared())
	})
	if err := generate.Literal(e, fn); err != nil {
		return err
	}
	e.Emit(isa.Call0)
	return nil
}
