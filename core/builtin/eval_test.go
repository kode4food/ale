package builtin_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/core/builtin"
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
	builtin.Eval(e1,
		add.Prepend(LS("list")),
	)
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

	// check to see that the third constant is evalFunc
	f, ok := c[2].(data.Procedure)
	as.True(ok)
	as.Equal(I(3), f.Call(add))
}
