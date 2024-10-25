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

func makeDefiner[T data.Value](m map[data.Local]T, err string) special.Call {
	return func(e encoder.Encoder, args ...data.Value) {
		data.AssertFixed(1, len(args))
		n := args[0].(data.Local)
		if sf, ok := m[n]; ok {
			if err := e.Globals().Declare(n).Bind(sf); err != nil {
				panic(err)
			}
			if err := generate.Symbol(e, n); err != nil {
				panic(err)
			}
			return
		}
		panic(fmt.Errorf(err, n))
	}
}
