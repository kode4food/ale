package bootstrap

import (
	"fmt"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/compiler/special"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

const (
	defBuiltInName = "def-builtin"
	defSpecialName = "def-special"
	defMacroName   = "def-macro"
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

func makeDefiner[T data.Value](m map[data.Local]T, errStr string) special.Call {
	return func(e encoder.Encoder, args ...data.Value) error {
		if err := data.CheckFixedArity(1, len(args)); err != nil {
			return err
		}
		n := args[0].(data.Local)
		if sf, ok := m[n]; ok {
			if err := env.BindPublic(e.Globals(), n, sf); err != nil {
				return err
			}
			return generate.Symbol(e, n)
		}
		return fmt.Errorf(errStr, n)
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
