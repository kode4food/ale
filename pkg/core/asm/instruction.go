package asm

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/strings"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/comb/basics"
)

func getInstructionCalls() namedAsmParsers {
	res := make(namedAsmParsers, len(isa.Effects))
	for oc, effect := range isa.Effects {
		name := data.Local(strings.CamelToSnake(oc.String()))
		res[name] = func(oc isa.Opcode, ao isa.ActOn) asmParse {
			return makeEmitCall(oc, ao)
		}(oc, effect.Operand)
	}
	return res
}

func makeEmitCall(oc isa.Opcode, actOn isa.ActOn) asmParse {
	if actOn == isa.Nothing {
		return makeStandaloneEmit(oc)
	}
	return makeOperandEmit(oc)
}

func makeStandaloneEmit(oc isa.Opcode) asmParse {
	return func(p *asmParser, s data.Sequence) (asmEmit, data.Sequence, error) {
		return func(e *asmEncoder) error {
			e.Emit(oc)
			return nil
		}, s, nil
	}
}

func makeOperandEmit(oc isa.Opcode) asmParse {
	return parseArgs(data.Local(oc.String()), 1,
		func(p *asmParser, args ...data.Value) (asmEmit, error) {
			return func(e *asmEncoder) error {
				e.Emit(oc, e.toOperands(oc, args)...)
				return nil
			}, nil
		},
	)
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
