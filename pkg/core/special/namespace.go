package special

import (
	"fmt"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/eval"
)

const (
	ErrExpectedName = "name expected, got %s"
)

func InNamespace(e encoder.Encoder, args ...data.Value) error {
	if err := data.CheckMinimumArity(2, len(args)); err != nil {
		return err
	}
	name, ok := args[0].(data.Local)
	if !ok {
		return fmt.Errorf(ErrExpectedName, args[0])
	}
	expr := data.Vector(args[1:])
	ns := e.Globals().Environment().GetQualified(name)
	fn := data.Call(func(...data.Value) data.Value {
		res, err := eval.Block(ns, expr)
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

func Import(e encoder.Encoder, args ...data.Value) error {
	if err := data.CheckRangedArity(1, 2, len(args)); err != nil {
		return err
	}
	name, ok := args[0].(data.Local)
	if !ok {
		return fmt.Errorf(ErrExpectedName, args[0])
	}
	from := e.Globals().Environment().GetQualified(name)
	to := e.Globals()
	fn, err := getImporter(from, to, args[1:]...)
	if err != nil {
		return err
	}
	if err := generate.Literal(e, fn); err != nil {
		return err
	}
	e.Emit(isa.Call0)
	return nil
}

func getImporter(from, to env.Namespace, args ...data.Value) (data.Call, error) {
	if len(args) == 0 {
		return importAll(from, to), nil
	}
	switch v := args[0].(type) {
	case data.Vector:
		return importNamed(from, to, v)
	case *data.Object:
		return importAliased(from, to, v)
	default:
		return nil, fmt.Errorf("nah, you can't do that")
	}
}

func importAll(from, to env.Namespace) data.Call {
	return func(...data.Value) data.Value {
		names := from.Declared()
		for _, n := range names {
			if err := copyEntry(from, to, n, n); err != nil {
				panic(err)
			}
		}
		return data.Null
	}
}

func importNamed(from, to env.Namespace, vals data.Vector) (data.Call, error) {
	names := make(data.Locals, len(vals))
	for i, v := range vals {
		if n, ok := v.(data.Local); ok {
			names[i] = n
			continue
		}
		return nil, fmt.Errorf(ErrExpectedName, v)
	}
	return func(...data.Value) data.Value {
		for _, n := range names {
			if err := copyEntry(from, to, n, n); err != nil {
				panic(err)
			}
		}
		return data.Null
	}, nil
}

func importAliased(from, to env.Namespace, a *data.Object) (data.Call, error) {
	aliases := map[data.Local]data.Local{}
	for _, p := range a.Pairs() {
		k, ok := p.Car().(data.Local)
		if !ok {
			return nil, fmt.Errorf(ErrExpectedName, p.Car())
		}
		v, ok := p.Cdr().(data.Local)
		if !ok {
			return nil, fmt.Errorf(ErrExpectedName, p.Cdr())
		}
		aliases[k] = v
	}
	return func(...data.Value) data.Value {
		for name, alias := range aliases {
			if err := copyEntry(from, to, name, alias); err != nil {
				panic(err)
			}
		}
		return data.Null
	}, nil
}

func copyEntry(from, to env.Namespace, name, alias data.Local) error {
	f, _, err := from.Resolve(name)
	if err != nil {
		return err
	}
	if !f.IsBound() {
		return nil
	}
	t, err := to.Public(alias)
	if err != nil {
		return err
	}
	v, err := f.Value()
	if err != nil {
		return err
	}
	return t.Bind(v)
}
