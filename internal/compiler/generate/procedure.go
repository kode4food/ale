package generate

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/procedure"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/runtime/vm"
)

type procEncoder struct {
	encoder.Encoder
}

func Procedure(e encoder.Encoder, build Builder) (*vm.Procedure, error) {
	child := makeProcEncoder(e).Child()
	if err := build(child); err != nil {
		return nil, err
	}
	enc := child.Encode()
	fn, err := procedure.FromEncoded(enc)
	if err != nil {
		return nil, err
	}

	cells := enc.Closure
	clen := len(cells)
	if clen == 0 {
		// nothing needed to be captured from local variables, so pass the
		// newly instantiated closure through
		if err := Literal(e, fn.Call()); err != nil {
			return nil, err
		}
		return fn, nil
	}

	for i := clen - 1; i >= 0; i-- {
		if err := Local(e, cells[i]); err != nil {
			return nil, err
		}
	}
	e.Emit(isa.Const, e.AddConstant(fn))
	e.Emit(isa.Call, isa.Operand(clen))
	return fn, nil
}

func makeProcEncoder(e encoder.Encoder) *procEncoder {
	return &procEncoder{
		Encoder: e,
	}
}

func (e *procEncoder) Wrapped() encoder.Encoder {
	return e.Encoder
}

func (e *procEncoder) Child() encoder.Encoder {
	res := *e
	res.Encoder = e.Encoder.Child()
	return &res
}
