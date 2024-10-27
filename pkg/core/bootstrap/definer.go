package bootstrap

import (
	"fmt"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/compiler/special"
	"github.com/kode4food/ale/pkg/data"
)

const (
	defBuiltInName = "def-builtin"
	defSpecialName = "def-special"
	defMacroName   = "def-macro"
)

func (b *bootstrap) populateDefiners() {
	ns := b.environment.GetRoot()

	_ = ns.Private(defBuiltInName).Bind(
		makeDefiner(b.procMap, ErrBuiltInNotFound),
	)
	_ = ns.Private(defSpecialName).Bind(
		makeDefiner(b.specialMap, ErrSpecialNotFound),
	)
	_ = ns.Private(defMacroName).Bind(
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
			if err := e.Globals().Declare(n).Bind(sf); err != nil {
				return err
			}
			return generate.Symbol(e, n)
		}
		return fmt.Errorf(errStr, n)
	}
}
