package asm

import (
	"errors"
	"fmt"
	"slices"

	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/strings"
	"github.com/kode4food/ale/pkg/data"
)

var asmEffects = excludeEffects([]isa.Opcode{
	isa.Call0, isa.Call1, isa.Call2, isa.Call3, isa.CallSelf, isa.Const,
	isa.TailCall, isa.TailClos, isa.TailSelf,
})

func getInstructionCalls() namedAsmParsers {
	res := make(namedAsmParsers, len(asmEffects))
	for oc, effect := range asmEffects {
		name := data.Local(strings.CamelToSnake(oc.String()))
		res[name] = func(oc isa.Opcode, ao isa.ActOn) asmParse {
			return makeEmitCall(oc, ao)
		}(oc, effect.Operand)
	}
	return res
}

func excludeEffects(exclude []isa.Opcode) map[isa.Opcode]*isa.Effect {
	res := make(map[isa.Opcode]*isa.Effect, len(isa.Effects)-len(exclude))
	for oc, effect := range isa.Effects {
		if slices.Contains(exclude, oc) {
			continue
		}
		res[oc] = effect
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
	return func(p *Parser, s data.Sequence) (Emit, data.Sequence, error) {
		return func(e *Encoder) error {
			e.Emit(oc)
			return nil
		}, s, nil
	}
}

func makeOperandEmit(oc isa.Opcode) asmParse {
	return parseArgs(data.Local(oc.String()), 1,
		func(p *Parser, args ...data.Value) (Emit, error) {
			return func(e *Encoder) error {
				ops, err := e.toOperands(oc, args)
				if err != nil {
					return err
				}
				e.Emit(oc, ops...)
				return nil
			}, nil
		},
	)
}

func (e *Encoder) toOperands(
	oc isa.Opcode, args data.Vector,
) ([]isa.Operand, error) {
	res := make([]isa.Operand, len(args))
	for i, a := range args {
		ao := isa.Effects[oc].Operand
		toOperand := e.getToOperandFor(ao)
		r, err := toOperand(e, a)
		if err != nil {
			return nil, err
		}
		res[i] = r
	}
	return res, nil
}

func (e *Encoder) getToOperandFor(ao isa.ActOn) asmToOperand {
	switch ao {
	case isa.Locals:
		return e.makeNameToWord()
	case isa.Labels:
		return e.makeLabelToWord()
	default:
		return toOperand
	}
}

func (e *Encoder) makeNameToWord() asmToOperand {
	return wrapToOperandError(ErrUnexpectedName,
		func(e *Encoder, val data.Value) (isa.Operand, error) {
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

func (e *Encoder) makeLabelToWord() asmToOperand {
	return wrapToOperandError(ErrUnexpectedLabel,
		func(e *Encoder, val data.Value) (isa.Operand, error) {
			if v, ok := e.resolveEncoderArg(val); ok {
				val = v
			}
			if val, ok := val.(data.Keyword); ok {
				return e.getLabelIndex(val), nil
			}
			return toOperand(e, val)
		},
	)
}

func wrapToOperandError(errStr string, toOperand asmToOperand) asmToOperand {
	return func(e *Encoder, val data.Value) (isa.Operand, error) {
		res, err := toOperand(e, val)
		if err != nil {
			return 0, errors.Join(fmt.Errorf(errStr, val), err)
		}
		return res, nil
	}
}

func toOperand(_ *Encoder, val data.Value) (isa.Operand, error) {
	if val, ok := val.(data.Integer); ok {
		if isa.IsValidOperand(int(val)) {
			return isa.Operand(val), nil
		}
	}
	return 0, fmt.Errorf(isa.ErrExpectedOperand, val)
}
