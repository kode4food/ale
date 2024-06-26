package generate

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/procedure"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/runtime/vm"
)

func Procedure(e encoder.Encoder, build Builder) *vm.Procedure {
	child := e.Child()
	build(child)
	enc := child.Encode()
	fn := procedure.FromEncoded(enc)

	cells := enc.Closure
	clen := len(cells)
	if clen == 0 {
		// nothing needed to be captured from local variables, so pass the
		// newly instantiated closure through
		Literal(e, fn.Call())
		return fn
	}

	for i := clen - 1; i >= 0; i-- {
		Symbol(e, cells[i])
	}
	e.Emit(isa.Const, e.AddConstant(fn))
	e.Emit(isa.Call, isa.Operand(clen))
	return fn
}
