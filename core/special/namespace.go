package special

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/runtime/isa"
)

type imports map[data.Local]data.Local

const (
	ErrExpectedName     = "name expected, got %s"
	ErrUnexpectedImport = "unexpected import pattern: %s"
	ErrDuplicateName    = "duplicate name(s) in import: %s"
)

func MakeNamespace(e encoder.Encoder, args ...ale.Value) error {
	if err := data.CheckFixedArity(1, len(args)); err != nil {
		return err
	}
	name, ok := args[0].(data.Local)
	if !ok {
		return fmt.Errorf(ErrExpectedName, args[0])
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
			return fmt.Errorf(ErrExpectedName, args[0])
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

func Import(e encoder.Encoder, args ...ale.Value) error {
	if err := data.CheckRangedArity(1, 2, len(args)); err != nil {
		return err
	}
	name, ok := args[0].(data.Local)
	if !ok {
		return fmt.Errorf(ErrExpectedName, args[0])
	}
	from, err := e.Globals().Environment().GetQualified(name)
	if err != nil {
		return err
	}
	to := e.Globals()
	fn, err := getImporter(from, to, args[1:]...)
	if err != nil {
		return err
	}
	proc := data.MakeProcedure(fn)
	if err := generate.Literal(e, proc); err != nil {
		return err
	}
	e.Emit(isa.Call0)
	return nil
}

func getImporter(from, to env.Namespace, args ...ale.Value) (data.Call, error) {
	if len(args) == 0 {
		return importAll(from, to), nil
	}
	switch v := args[0].(type) {
	case data.Local:
		return importNamed(from, to, data.NewList(v))
	case data.Vector:
		return importNamed(from, to, data.NewList(v))
	case *data.List:
		return importNamed(from, to, v)
	default:
		return nil, fmt.Errorf(ErrUnexpectedImport, args[0])
	}
}

func importAll(from, to env.Namespace) data.Call {
	return func(...ale.Value) ale.Value {
		names := localsToVector(from.Declared())
		i, err := buildImports(data.NewList(names...))
		if err != nil {
			panic(err)
		}
		if err := performImports(from, to, i); err != nil {
			panic(err)
		}
		return localsToVector(basics.MapValues(i))
	}
}

func importNamed(from, to env.Namespace, a *data.List) (data.Call, error) {
	i, err := buildImports(a)
	if err != nil {
		return nil, err
	}
	return func(...ale.Value) ale.Value {
		err := performImports(from, to, i)
		if err != nil {
			panic(err)
		}
		return localsToVector(basics.MapValues(i))
	}, nil
}

func buildImports(a *data.List) (imports, error) {
	res := imports{}
	for f, r, ok := a.Split(); ok; f, r, ok = r.Split() {
		switch f := f.(type) {
		case data.Local:
			if _, ok := res[f]; ok {
				return nil, fmt.Errorf(ErrDuplicateName, f)
			}
			res[f] = f
		case data.Vector:
			if len(f) != 2 {
				return nil, errors.New(ErrUnpairedBindings)
			}
			alias, ok := f[0].(data.Local)
			if !ok {
				return nil, fmt.Errorf(ErrExpectedName, f[0])
			}
			name, ok := f[1].(data.Local)
			if !ok {
				return nil, fmt.Errorf(ErrExpectedName, f[1])
			}
			if _, ok := res[alias]; ok {
				return nil, fmt.Errorf(ErrDuplicateName, alias)
			}
			res[alias] = name
		default:
			return nil, fmt.Errorf(ErrUnexpectedImport, f)
		}
	}
	return res, nil
}

func performImports(from, to env.Namespace, i imports) error {
	le := map[data.Local]*env.Entry{}
	for alias, name := range i {
		e, _, err := from.Resolve(name)
		if err != nil {
			return err
		}
		if e.IsPrivate() {
			return fmt.Errorf(env.ErrNameNotDeclared, name)
		}
		le[alias] = e
	}
	return to.Import(le)
}

func localsToVector(locals data.Locals) data.Vector {
	res := make(data.Vector, len(locals))
	for i, n := range locals {
		res[i] = n
	}
	return res
}
