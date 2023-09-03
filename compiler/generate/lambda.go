package generate

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/ale/runtime/vm"
)

func Lambda(e encoder.Encoder, build Builder) *vm.Lambda {
	child := e.Child()
	build(child)
	fn := vm.LambdaFromEncoder(child)

	cells := child.Closure()
	clen := len(cells)
	if clen == 0 {
		// nothing needed to be captured from local variables,
		// so just pass the newly instantiated closure through
		Literal(e, fn.Call())
		return fn
	}

	for i := clen - 1; i >= 0; i-- {
		name := cells[i].Name
		Symbol(e, data.NewLocalSymbol(name))
	}
	e.Emit(isa.Const, e.AddConstant(fn))
	e.Emit(isa.Call, isa.Operand(clen))
	return fn
}
