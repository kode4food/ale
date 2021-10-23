package generate_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
)

func TestPair(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	generate.Value(e1, data.NewCons(S("left"), S("right")))
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.New(isa.Const, 0),
		isa.New(isa.Const, 1),
		isa.New(isa.Const, 2),
		isa.New(isa.Call, 2),
		isa.New(isa.Return),
	}, e1.Code())

	c := e1.Constants()
	as.Equal(S("right"), c[0])
	as.Equal(S("left"), c[1])

	cons := assert.GetRootSymbol(e1, "cons")
	as.Equal(cons, c[2])
}
