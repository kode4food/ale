package asm

import (
	"fmt"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/comb/basics"
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
	// ErrBadNameResolution is raised when an attempt is made to bind a local
	// using an argument to the encoder that is not a Local symbol
	ErrBadNameResolution = "encoder argument is not a name: %s"
)

var gen = data.NewSymbolGenerator()

func noAsmEmit(*asmEncoder) error { return nil }

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

func (e *asmEncoder) toOperands(oc isa.Opcode, args data.Vector) []isa.Operand {
	return basics.Map(args, func(a data.Value) isa.Operand {
		ao := isa.Effects[oc].Operand
		toOperand := e.getToOperandFor(ao)
		r, err := toOperand(e, a)
		if err != nil {
			panic(err)
		}
		return r
	})
}

func (e *asmEncoder) getToOperandFor(ao isa.ActOn) asmToOperand {
	switch ao {
	case isa.Locals:
		return e.makeNameToWord()
	case isa.Labels:
		return e.makeLabelToWord()
	default:
		return toOperand
	}
}

func (e *asmEncoder) makeLabelToWord() asmToOperand {
	return wrapToOperandError(ErrUnexpectedLabel,
		func(e *asmEncoder, val data.Value) (isa.Operand, error) {
			if v, ok := e.resolveEncoderArg(val); ok {
				val = v
			}
			if val, ok := val.(data.Keyword); ok {
				return e.getLabelIndex(val.Name()), nil
			}
			return toOperand(e, val)
		},
	)
}

func (e *asmEncoder) makeNameToWord() asmToOperand {
	return wrapToOperandError(ErrUnexpectedName,
		func(e *asmEncoder, val data.Value) (isa.Operand, error) {
			if v, ok := e.resolveEncoderArg(val); ok {
				val = v
			}
			if val, ok := val.(data.Local); ok {
				n := e.resolvePrivate(val)
				if cell, ok := e.ResolveLocal(n); ok {
					return cell.Index, nil
				}
				return 0, fmt.Errorf(ErrUnexpectedName, val)
			}
			return toOperand(e, val)
		},
	)
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
