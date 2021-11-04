package generate_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
)

func TestBlock(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	generate.Block(e1, V())
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.New(isa.Nil),
		isa.New(isa.Return),
	}, e1.Code())

	e2 := assert.GetTestEncoder()
	generate.Block(e2, V(
		L(LS("+"), I(1), I(2)),
		B(true),
	))
	e2.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.New(isa.Two),
		isa.New(isa.One),
		isa.New(isa.Const, 0),
		isa.New(isa.Call, 2),
		isa.New(isa.Pop),
		isa.New(isa.True),
		isa.New(isa.Return),
	}, e2.Code())

	c := e2.Constants()
	as.Equal(assert.GetRootSymbol(e2, "+"), c[0])
}
