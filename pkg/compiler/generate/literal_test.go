package generate_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/compiler/generate"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/runtime/isa"
)

func TestLiteral(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	generate.Literal(e, I(0))
	generate.Literal(e, I(1))
	generate.Literal(e, I(2))
	generate.Literal(e, I(3))
	generate.Literal(e, I(-1))
	generate.Literal(e, data.True)
	generate.Literal(e, data.False)
	generate.Literal(e, data.Null)
	generate.Literal(e, S("hello there!"))

	// Because the stack size must remain the same in and out
	for i := 0; i < 9; i++ {
		e.Emit(isa.Pop)
	}

	enc := e.Encode()
	as.Instructions(
		isa.Instructions{
			isa.Zero.New(),
			isa.PosInt.New(1),
			isa.PosInt.New(2),
			isa.PosInt.New(3),
			isa.NegInt.New(1),
			isa.True.New(),
			isa.False.New(),
			isa.Null.New(),
			isa.Const.New(0),
			isa.Pop.New(),
			isa.Pop.New(),
			isa.Pop.New(),
			isa.Pop.New(),
			isa.Pop.New(),
			isa.Pop.New(),
			isa.Pop.New(),
			isa.Pop.New(),
			isa.Pop.New(),
		},
		enc.Code,
	)

	c := enc.Constants
	as.Equal(S("hello there!"), c[0])
}
