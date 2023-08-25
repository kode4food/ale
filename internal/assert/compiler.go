package assert

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/runtime/isa"
)

// Instructions tests that two sets of Instructions are identical
func (w *Wrapper) Instructions(expected, actual isa.Instructions) {
	w.Helper()
	w.Equal(len(expected), len(actual))
	for i, l := range expected {
		a := actual[i]
		w.Assertions.Equal(l.Opcode, a.Opcode)
		w.Assertions.Equal(len(l.Args), len(a.Args))
		if len(l.Args) > 0 {
			w.Assertions.Equal(l.Args, a.Args)
		}
	}
}

// GetRootSymbol is a test helper that retrieves the value for a named symbol
// from the Encoder's global environment or dies trying
func GetRootSymbol(e encoder.Encoder, n data.Name) data.Value {
	s := env.RootSymbol(n)
	ge := e.Globals().Environment()
	root := ge.GetRoot()
	return env.MustResolveValue(root, s)
}
