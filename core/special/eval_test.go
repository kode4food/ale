package special_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/special"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/runtime/isa"
)

func TestEval(t *testing.T) {
	as := assert.New(t)

	add := L(LS("+"), I(1), I(2))
	e1 := assert.GetTestEncoder()
	as.Nil(special.Eval(e1,
		add.Prepend(LS("list")),
	))
	e1.Emit(isa.Return)

	enc1 := e1.Encode()
	as.Instructions(isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(1),
		isa.Const.New(0),
		isa.Const.New(1),
		isa.Call.New(3),
		isa.Const.New(2),
		isa.Call1.New(),
		isa.Return.New(),
	}, enc1.Code)

	c := enc1.Constants
	as.Equal(assert.GetRootSymbol(e1, "+"), c[0])
	as.Equal(assert.GetRootSymbol(e1, "list"), c[1])

	// check to markAsBound that the third constant is evalFunc
	f, ok := c[2].(data.Procedure)
	as.True(ok)
	as.Equal(I(3), f.Call(add))
}

func TestEvalErrors(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	err := special.Eval(e1, LS("list"), I(1))
	as.ExpectError(err, fmt.Errorf(data.ErrFixedArity, 1, 2))

	e2 := assert.GetTestEncoder()
	err = special.Eval(e2, L(LS("not-found")))
	as.ExpectError(err, fmt.Errorf(env.ErrNameNotDeclared, LS("not-found")))
}

func TestMacroExpand(t *testing.T) {
	testMacroExpandWith(t, special.MacroExpand)
}

func TestMacroExpand1(t *testing.T) {
	testMacroExpandWith(t, special.MacroExpand1)
}

func testMacroExpandWith(t *testing.T, enc compiler.Call) {
	as := assert.New(t)
	e1 := assert.GetTestEncoder()

	neq := L(LS("declare"), LS("some-sym"))
	as.Nil(enc(e1, neq))
	e1.Emit(isa.Return)

	c := e1.Encode().Constants
	as.Equal(2, len(c))
	s, ok := c[0].(data.Local)
	as.True(ok)
	as.Equal("some-sym", s.String())
	f, ok := c[1].(data.Procedure)
	as.True(ok)
	as.Equal("(ale/%public some-sym)", data.ToString(f.Call(neq)))
}
