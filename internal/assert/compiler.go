package assert

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/read"
	"github.com/kode4food/ale/runtime/isa"
)

// EncodesAs tests that a string generates the expected set of Instructions
func (w *Wrapper) EncodesAs(expected isa.Instructions, src data.String) {
	e := GetTestEncoder()
	v := read.FromString(src)
	generate.Block(e, v)
	w.Instructions(expected, e.Code())
}

// Instructions test that two sets of Instructions are identical
func (w *Wrapper) Instructions(expected, actual isa.Instructions) {
	w.Helper()
	w.Equal(len(expected), len(actual))
	for i, e := range expected {
		w.Assertions.Equal(e, actual[i])
	}
}

// GetRootSymbol is a test helper that retrieves the value for a named symbol
// from the Encoder's global environment or dies trying
func GetRootSymbol(e encoder.Encoder, n data.Local) data.Value {
	s := env.RootSymbol(n)
	ge := e.Globals().Environment()
	root := ge.GetRoot()
	return env.MustResolveValue(root, s)
}
