package generate_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
)

func TestLiteral(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	generate.Literal(e, data.Integer(0))
	generate.Literal(e, data.Integer(1))
	generate.Literal(e, data.Integer(2))
	generate.Literal(e, data.Integer(3))
	generate.Literal(e, data.Integer(-1))
	generate.Literal(e, data.True)
	generate.Literal(e, data.False)
	generate.Literal(e, data.Nil)
	generate.Literal(e, S("hello there!"))

	// Because the stack size must remain the same in and out
	for i := 0; i < 9; i++ {
		e.Emit(isa.Pop)
	}

	as.Instructions(
		isa.Instructions{
			isa.New(isa.Zero),
			isa.New(isa.One),
			isa.New(isa.Two),
			isa.New(isa.Const, 0),
			isa.New(isa.NegOne),
			isa.New(isa.True),
			isa.New(isa.False),
			isa.New(isa.Nil),
			isa.New(isa.Const, 1),
			isa.New(isa.Pop),
			isa.New(isa.Pop),
			isa.New(isa.Pop),
			isa.New(isa.Pop),
			isa.New(isa.Pop),
			isa.New(isa.Pop),
			isa.New(isa.Pop),
			isa.New(isa.Pop),
			isa.New(isa.Pop),
		},
		e.Code(),
	)
}
