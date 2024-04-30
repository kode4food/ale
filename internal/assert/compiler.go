package assert

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/read"
)

// EncodesAs tests that a string generates the expected set of Instructions
func (w *Wrapper) EncodesAs(expected isa.Instructions, src data.String) {
	e := GetTestEncoder()
	v := read.FromString(src)
	generate.Block(e, v)
	w.Instructions(expected, e.Encode().Code)
}

// Instructions test that two sets of Instructions are identical
func (w *Wrapper) Instructions(expected, actual isa.Instructions) {
	w.Helper()
	w.Equal(expected.String(), actual.String())
}

// GetRootSymbol is a test helper that retrieves the value for a named symbol
// from the Encoder's global environment or dies trying
func GetRootSymbol(e encoder.Encoder, n data.Local) data.Value {
	s := env.RootSymbol(n)
	ge := e.Globals().Environment()
	root := ge.GetRoot()
	return env.MustResolveValue(root, s)
}
