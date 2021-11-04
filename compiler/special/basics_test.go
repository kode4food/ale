package special_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
)

func TestBegin(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	special.Begin(e1,
		L(LS("+"), I(1), I(2)),
		B(true),
	)
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.New(isa.Two),
		isa.New(isa.One),
		isa.New(isa.Const, 0),
		isa.New(isa.Call, 2),
		isa.New(isa.Pop),
		isa.New(isa.True),
		isa.New(isa.Return),
	}, e1.Code())

	c := e1.Constants()
	as.Equal(assert.GetRootSymbol(e1, "+"), c[0])
}

func TestEval(t *testing.T) {
	as := assert.New(t)

	add := L(LS("+"), I(1), I(2))
	e1 := assert.GetTestEncoder()
	special.Eval(e1,
		add.Prepend(LS("list")),
	)
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.New(isa.Two),
		isa.New(isa.One),
		isa.New(isa.Const, 0),
		isa.New(isa.Resolve),
		isa.New(isa.Const, 1),
		isa.New(isa.Call, 3),
		isa.New(isa.Const, 2),
		isa.New(isa.Call1),
		isa.New(isa.Return),
	}, e1.Code())

	c := e1.Constants()
	as.Equal(LS("+"), c[0])
	as.Equal(assert.GetRootSymbol(e1, "list"), c[1])

	// check to see that the third constant is evalFor
	f, ok := c[2].(data.Function)
	as.True(ok)
	as.Equal(I(3), f.Call(add))
}
