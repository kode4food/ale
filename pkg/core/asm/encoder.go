package asm

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
)

type (
	asmEncoder struct {
		encoder.Encoder
		*asmParser
		args    map[data.Local]data.Value
		labels  map[data.Local]isa.Operand
		private map[data.Local]data.Local
	}

	asmEmit      func(*asmEncoder) error
	asmToOperand func(*asmEncoder, data.Value) (isa.Operand, error)
	asmToName    func(*asmEncoder, data.Local) (data.Local, error)
)

const (
	// ErrUnexpectedName is raised when a local name is referenced that hasn't
	// been declared as part of the assembler encoder's scope
	ErrUnexpectedName = "unexpected local name: %s"

	// ErrUnexpectedLabel is raised when a jump or cond-jump instruction refers
	// to a label that hasn't been anchored in the assembler block
	ErrUnexpectedLabel = "unexpected label: %s"
)

var gen = data.NewSymbolGenerator()

func noAsmEmit(_ *asmEncoder) error { return nil }

func (p *asmParser) wrapEncoder(
	e encoder.Encoder, args ...data.Value,
) *asmEncoder {
	a := make(map[data.Local]data.Value, len(args))
	for i, k := range p.params {
		a[k] = args[i]
	}
	return &asmEncoder{
		asmParser: p,
		Encoder:   e,
		args:      a,
		labels:    map[data.Local]isa.Operand{},
		private:   map[data.Local]data.Local{},
	}
}

func (e *asmEncoder) Wrapped() encoder.Encoder {
	return e.Encoder
}

func (e *asmEncoder) Child() encoder.Encoder {
	res := *e
	res.Encoder = e.Encoder.Child()
	return &res
}

func (e *asmEncoder) getLabelIndex(n data.Local) isa.Operand {
	if idx, ok := e.labels[n]; ok {
		return idx
	}
	idx := e.NewLabel()
	e.labels[n] = idx
	return idx
}

func (e *asmEncoder) resolvePrivate(l data.Local) data.Local {
	if g, ok := e.private[l]; ok {
		return g
	}
	return l
}

func (e *asmEncoder) resolveEncoderArg(v data.Value) (data.Value, bool) {
	if v, ok := v.(data.Local); ok {
		if res, ok := e.args[v]; ok {
			return res, true
		}
	}
	if p, ok := e.Encoder.(*asmEncoder); ok {
		return p.resolveEncoderArg(v)
	}
	return nil, false
}
