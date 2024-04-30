package generate_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/runtime/isa"
)

func TestBlock(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	generate.Block(e1, V())
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.Null.New(),
		isa.Return.New(),
	}, e1.Encode().Code)

	e2 := assert.GetTestEncoder()
	generate.Block(e2, V(
		L(LS("+"), I(1), I(2)),
		B(true),
	))
	e2.Emit(isa.Return)

	enc2 := e2.Encode()
	as.Instructions(isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(1),
		isa.Const.New(0),
		isa.Call.New(2),
		isa.Pop.New(),
		isa.True.New(),
		isa.Return.New(),
	}, enc2.Code)

	c := enc2.Constants
	as.Equal(assert.GetRootSymbol(e2, "+"), c[0])
}
