package generate_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
)

func TestLiteral(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	as.Nil(generate.Literal(e, I(0)))
	as.Nil(generate.Literal(e, I(1)))
	as.Nil(generate.Literal(e, I(2)))
	as.Nil(generate.Literal(e, I(3)))
	as.Nil(generate.Literal(e, I(-1)))
	as.Nil(generate.Literal(e, data.True))
	as.Nil(generate.Literal(e, data.False))
	as.Nil(generate.Literal(e, data.Null))
	as.Nil(generate.Literal(e, S("hello there!")))

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
