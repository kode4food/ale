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

	ns.Private(defBuiltInName).Bind(
		makeDefiner(b.procMap, ErrBuiltInNotFound),
	)
	ns.Private(defSpecialName).Bind(
		makeDefiner(b.specialMap, ErrSpecialNotFound),
	)
	ns.Private(defMacroName).Bind(
		makeDefiner(b.macroMap, ErrMacroNotFound),
	)
}

func makeDefiner[T data.Value](m map[data.Local]T, err string) special.Call {
	return func(e encoder.Encoder, args ...data.Value) {
		data.AssertFixed(1, len(args))
		n := args[0].(data.Local)
		if sf, ok := m[n]; ok {
			e.Globals().Declare(n).Bind(sf)
			generate.Symbol(e, n)
			return
		}
		panic(fmt.Errorf(err, n))
	}
}
