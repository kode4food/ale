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

func TestEval(t *testing.T) {
	as := assert.New(t)

	add := L(LS("+"), I(1), I(2))
	e1 := assert.GetTestEncoder()
	special.Eval(e1,
		add.Prepend(LS("list")),
	)
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(1),
		isa.Const.New(0),
		isa.Resolve.New(),
		isa.Const.New(1),
		isa.Call.New(3),
		isa.Const.New(2),
		isa.Call1.New(),
		isa.Return.New(),
	}, e1.Code())

	c := e1.Constants()
	as.Equal(LS("+"), c[0])
	as.Equal(assert.GetRootSymbol(e1, "list"), c[1])

	// check to see that the third constant is evalFunc
	f, ok := c[2].(data.Function)
	as.True(ok)
	as.Equal(I(3), f.Call(add))
}

func TestMacroExpand(t *testing.T) {
	testMacroExpandWith(t, special.MacroExpand)
}

func TestMacroExpand1(t *testing.T) {
	testMacroExpandWith(t, special.MacroExpand1)
}

func testMacroExpandWith(t *testing.T, enc testEncoder) {
	as := assert.New(t)
	e1 := assert.GetTestEncoder()

	neq := L(LS("declare"), LS("some-sym"))
	enc(e1, neq)
	e1.Emit(isa.Return)

	c := e1.Constants()
	as.Equal(2, len(c))
	s, ok := c[0].(data.Local)
	as.True(ok)
	as.Equal("some-sym", s.String())
	f, ok := c[1].(data.Function)
	as.True(ok)
	as.Equal("(ale/declare* some-sym)", f.Call(neq).String())
}

func TestBegin(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	special.Begin(e1,
		L(LS("+"), I(1), I(2)),
		B(true),
	)
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(1),
		isa.Const.New(0),
		isa.Call.New(2),
		isa.Pop.New(),
		isa.True.New(),
		isa.Return.New(),
	}, e1.Code())

	c := e1.Constants()
	as.Equal(assert.GetRootSymbol(e1, "+"), c[0])
}
