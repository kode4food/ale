package bootstrap

import (
	"fmt"

	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

const (
	defBuiltInName = "def-builtin"
	defSpecialName = "def-special"
	defMacroName   = "def-macro"

	kwdPublic  = data.Keyword("public")
	kwdPrivate = data.Keyword("private")
)

func (b *bootstrap) populateDefiners() {
	ns := b.environment.GetRoot()

	mustBindPrivate(ns, defBuiltInName,
		makeDefiner(b.procMap, ErrBuiltInNotFound),
	)
	mustBindPrivate(ns, defSpecialName,
		makeDefiner(b.specialMap, ErrSpecialNotFound),
	)
	mustBindPrivate(ns, defMacroName,
		makeDefiner(b.macroMap, ErrMacroNotFound),
	)
}

func makeDefiner[T data.Value](
	m map[data.Local]T, errStr string,
) compiler.Call {
	return func(e encoder.Encoder, args ...data.Value) error {
		if err := data.CheckRangedArity(1, 2, len(args)); err != nil {
			return err
		}
		bind, args, err := getBinder(args...)
		if err != nil {
			return err
		}
		n := args[0].(data.Local)
		if sf, ok := m[n]; ok {
			if err := bind(e.Globals(), n, sf); err != nil {
				return err
			}
			return generate.Local(e, n)
		}
		return fmt.Errorf(errStr, n)
	}
}

func getBinder(args ...data.Value) (env.Binder, []data.Value, error) {
	k, ok := args[0].(data.Keyword)
	if !ok {
		return env.BindPublic, args, nil
	}
	switch k {
	case kwdPublic:
		return env.BindPublic, args[1:], nil
	case kwdPrivate:
		return env.BindPrivate, args[1:], nil
	default:
		return nil, nil, fmt.Errorf("unknown binding keyword: %s", k)
	}
}

func mustBindPublic(ns env.Namespace, n data.Local, v data.Value) {
	if err := env.BindPublic(ns, n, v); err != nil {
		panic(err)
	}
}

func mustBindPrivate(ns env.Namespace, n data.Local, v data.Value) {
	if err := env.BindPrivate(ns, n, v); err != nil {
		panic(err)
	}
}
