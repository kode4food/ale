package asm

import (
	"errors"
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
	// ErrUnexpectedName is raised when a local name is referenced that hasn't
	// been declared as part of the assembler encoder's scope
	ErrUnexpectedName = "unexpected local name: %s"

	// ErrUnexpectedLabel is raised when a jump or cond-jump instruction refers
	// to a label that hasn't been anchored in the assembler block
	ErrUnexpectedLabel = "unexpected label: %s"
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

func wrapToOperandError(errStr string, toOperand asmToOperand) asmToOperand {
	return func(e *asmEncoder, val data.Value) (isa.Operand, error) {
		res, err := toOperand(e, val)
		if err != nil {
			return 0, errors.Join(fmt.Errorf(errStr, val), err)
		}
		return res, nil
	}
}

func toOperand(_ *asmEncoder, val data.Value) (isa.Operand, error) {
	if val, ok := val.(data.Integer); ok {
		if isa.IsValidOperand(int(val)) {
			return isa.Operand(val), nil
		}
	}
	return 0, fmt.Errorf(isa.ErrExpectedOperand, val)
}
