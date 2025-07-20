package asm

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/runtime/isa"
)

type (
	Encoder struct {
		encoder.Encoder
		*Parser
		args    map[data.Local]ale.Value
		labels  map[data.Keyword]isa.Operand
		private map[data.Local]data.Local
	}

	asmToOperand func(*Encoder, ale.Value) (isa.Operand, error)
	asmToName    func(*Encoder, data.Local) (data.Local, error)
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

func noAsmEmit(_ *Encoder) error { return nil }

func (p *Parser) wrapEncoder(
	e encoder.Encoder, args ...ale.Value,
) *Encoder {
	a := make(map[data.Local]ale.Value, len(args))
	for i, k := range p.params {
		a[k] = args[i]
	}
	return &Encoder{
		Parser:  p,
		Encoder: e,
		args:    a,
		labels:  map[data.Keyword]isa.Operand{},
		private: map[data.Local]data.Local{},
	}
}

func (e *Encoder) Wrapped() encoder.Encoder {
	return e.Encoder
}

func (e *Encoder) Child() encoder.Encoder {
	res := *e
	res.Encoder = e.Encoder.Child()
	return &res
}

func (e *Encoder) getLabelIndex(n data.Keyword) isa.Operand {
	if idx, ok := e.labels[n]; ok {
		return idx
	}
	idx := e.NewLabel()
	e.labels[n] = idx
	return idx
}

func (e *Encoder) resolvePrivate(l data.Local) data.Local {
	if g, ok := e.private[l]; ok {
		return g
	}
	return l
}

func (e *Encoder) resolveEncoderArg(v ale.Value) (ale.Value, bool) {
	if v, ok := v.(data.Local); ok {
		if res, ok := e.args[v]; ok {
			return res, true
		}
	}
	if p, ok := e.Encoder.(*Encoder); ok {
		return p.resolveEncoderArg(v)
	}
	return nil, false
}
