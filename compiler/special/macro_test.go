package special_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
)

type testEncoder func(encoder.Encoder, ...data.Value)

func TestMacroExpand(t *testing.T) {
	testMacroExpandWith(t, special.MacroExpand)
}

func TestMacroExpand1(t *testing.T) {
	testMacroExpandWith(t, special.MacroExpand1)
}

func testMacroExpandWith(t *testing.T, enc testEncoder) {
	as := assert.New(t)
	e1 := assert.GetTestEncoder()

	neq := L(LS("!eq"), I(1), I(2))
	enc(e1, neq)
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.New(isa.Two),
		isa.New(isa.One),
		isa.New(isa.Const, 0),
		isa.New(isa.Call, 2),
		isa.New(isa.Const, 1),
		isa.New(isa.Call1),
		isa.New(isa.Const, 2),
		isa.New(isa.Call1),
		isa.New(isa.Return),
	}, e1.Code())

	c := e1.Constants()
	as.Equal(assert.GetRootSymbol(e1, "eq"), c[0])
	as.Equal(assert.GetRootSymbol(e1, "not"), c[1])

	// check to see that the third constant is expandFor
	f, ok := c[2].(data.Function)
	as.True(ok)
	as.Equal("(not (ale/eq 1 2))", f.Call(neq).String())
}
