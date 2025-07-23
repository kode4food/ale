package bootstrap

import (
	"fmt"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
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

func makeDefiner[T ale.Value](
	m map[data.Local]T, errStr string,
) compiler.Call {
	return func(e encoder.Encoder, args ...ale.Value) error {
		if err := data.CheckFixedArity(1, len(args)); err != nil {
			return err
		}
		n := args[0].(data.Local)
		if sf, ok := m[n]; ok {
			if err := env.BindPublic(e.Globals(), n, sf); err != nil {
				return err
			}
			return generate.Local(e, n)
		}
		return fmt.Errorf(errStr, n)
	}
}

func mustBindPublic(ns env.Namespace, n data.Local, v ale.Value) {
	if err := env.BindPublic(ns, n, v); err != nil {
		panic(err)
	}
}

func mustBindPrivate(ns env.Namespace, n data.Local, v ale.Value) {
	if err := env.BindPrivate(ns, n, v); err != nil {
		panic(err)
	}
}
