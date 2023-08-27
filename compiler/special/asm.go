package special

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/strings"
	"github.com/kode4food/ale/runtime/isa"
)

type (
	asmEncoder struct {
		encoder.Encoder
		labels map[data.Name]isa.Index
	}

	call struct {
		encoder.Call
		argCount int
	}

	callMap map[data.Name]*call

	toWordFunc func(data.Value) (isa.Word, error)
)

// Error messages
const (
	ErrUnknownDirective      = "unknown directive: %s"
	ErrUnexpectedForm        = "unexpected form: %s"
	ErrIncompleteInstruction = "incomplete instruction: %s"
	ErrUnknownLocalType      = "unknown local type: %s"
	ErrUnexpectedName        = "unexpected local name: %s"
	ErrUnexpectedLabel       = "unexpected label: %s"
	ErrExpectedWord          = "expected unsigned word: %s"
)

const (
	Value      = data.Name(".value")
	Const      = data.Name(".const")
	Local      = data.Name(".local")
	PushLocals = data.Name(".push-locals")
	PopLocals  = data.Name(".pop-locals")
)

var (
	instructionCalls = getInstructionCalls()
	encoderCalls     = getEncoderCalls()
	calls            = mergeCalls(instructionCalls, encoderCalls)

	cellTypes = map[data.Keyword]encoder.CellType{
		data.Keyword("val"):  encoder.ValueCell,
		data.Keyword("ref"):  encoder.ReferenceCell,
		data.Keyword("rest"): encoder.RestCell,
	}
)

// Asm provides indirect access to the Encoder's methods and generators
func Asm(e encoder.Encoder, args ...data.Value) {
	ae := &asmEncoder{
		Encoder: e,
		labels:  map[data.Name]isa.Index{},
	}
	v := data.NewVector(args...)
	for f, r, ok := v.Split(); ok; f, r, ok = r.Split() {
		switch v := f.(type) {
		case data.Keyword:
			e.Emit(isa.Label, ae.getLabelIndex(v))
		case data.LocalSymbol:
			n := v.Name()
			if d, ok := calls[n]; ok {
				if args, rest, ok := take(r, d.argCount); ok {
					d.Call(ae, args...)
					r = rest
					continue
				}
				panic(fmt.Errorf(ErrIncompleteInstruction, n))
			}
			panic(fmt.Errorf(ErrUnknownDirective, n))
		default:
			panic(fmt.Errorf(ErrUnexpectedForm, f.String()))
		}
	}
}

func getInstructionCalls() callMap {
	res := make(callMap, len(isa.Effects))
	for oc, effect := range isa.Effects {
		name := data.Name(strings.CamelToSnake(oc.String()))
		res[name] = func(oc isa.Opcode, argCount int) *call {
			return makeEmitCall(oc, argCount)
		}(oc, effect.Size-1)
	}
	return res
}

func makeEmitCall(oc isa.Opcode, argCount int) *call {
	return &call{
		Call: func(e encoder.Encoder, args ...data.Value) {
			e.Emit(oc, e.(*asmEncoder).toWords(oc, args)...)
		},
		argCount: argCount,
	}
}

func getEncoderCalls() callMap {
	return callMap{
		Value: {
			Call: func(e encoder.Encoder, args ...data.Value) {
				generate.Value(e, args[0])
			},
			argCount: 1,
		},
		Const: {
			Call: func(e encoder.Encoder, args ...data.Value) {
				index := e.AddConstant(args[0])
				e.Emit(isa.Const, index)
			},
			argCount: 1,
		},
		Local: {
			Call: func(e encoder.Encoder, args ...data.Value) {
				name := args[0].(data.LocalSymbol).Name()
				kwd := args[1].(data.Keyword)
				cellType, ok := cellTypes[kwd]
				if !ok {
					panic(fmt.Errorf(ErrUnknownLocalType, kwd))
				}
				e.AddLocal(name, cellType)
			},
			argCount: 2,
		},
		PushLocals: {
			Call: func(e encoder.Encoder, _ ...data.Value) {
				e.PushLocals()
			},
		},
		PopLocals: {
			Call: func(e encoder.Encoder, _ ...data.Value) {
				e.PopLocals()
			},
		},
	}
}

func mergeCalls(maps ...callMap) callMap {
	res := callMap{}
	for _, m := range maps {
		for k, v := range m {
			if _, ok := res[k]; ok {
				// Programmer error
				panic(fmt.Sprintf("duplicate entry: %s", k))
			}
			res[k] = v
		}
	}
	return res
}

func take(s data.Sequence, count int) (data.Values, data.Sequence, bool) {
	var f data.Value
	var ok bool
	res := make(data.Values, count)
	for i := 0; i < count; i++ {
		if f, s, ok = s.Split(); !ok {
			return nil, nil, false
		}
		res[i] = f
	}
	return res, s, true
}

func (e *asmEncoder) getLabelIndex(k data.Keyword) isa.Index {
	name := k.Name()
	if idx, ok := e.labels[name]; ok {
		return idx
	}
	idx := e.NewLabel()
	e.labels[name] = idx
	return idx
}

func (e *asmEncoder) toWords(oc isa.Opcode, args data.Values) []isa.Coder {
	toWord := e.getToWordFor(isa.Effects[oc])
	res := make([]isa.Coder, len(args))
	for i, a := range args {
		r, err := toWord(a)
		if err != nil {
			panic(err)
		}
		res[i] = r
	}
	return res
}

func (e *asmEncoder) getToWordFor(effect *isa.Effect) toWordFunc {
	fn := toWord
	if effect.Locals {
		fn = makeNameToWord(e, fn)
	}
	if effect.Labels {
		fn = makeLabelToWord(e, fn)
	}
	return fn
}

func makeLabelToWord(e *asmEncoder, next toWordFunc) toWordFunc {
	return wrapToWordError(func(val data.Value) (isa.Word, error) {
		if val, ok := val.(data.Keyword); ok {
			return isa.Word(e.getLabelIndex(val)), nil
		}
		return next(val)
	}, ErrUnexpectedLabel)
}

func makeNameToWord(e *asmEncoder, next toWordFunc) toWordFunc {
	return wrapToWordError(func(val data.Value) (isa.Word, error) {
		if val, ok := val.(data.LocalSymbol); ok {
			if cell, ok := e.ResolveLocal(val.Name()); ok {
				return isa.Word(cell.Index), nil
			}
			return 0, fmt.Errorf(ErrUnexpectedName, val)
		}
		return next(val)
	}, ErrUnexpectedName)
}

func wrapToWordError(toWord toWordFunc, errStr string) toWordFunc {
	return func(val data.Value) (isa.Word, error) {
		res, err := toWord(val)
		if err != nil {
			return 0, errors.Join(fmt.Errorf(errStr, val), err)
		}
		return res, nil
	}
}

func toWord(val data.Value) (isa.Word, error) {
	if val, ok := val.(data.Integer); ok {
		if isValidWord(val) {
			return isa.Word(val), nil
		}
	}
	return 0, fmt.Errorf(ErrExpectedWord, val)
}

func isValidWord(i data.Integer) bool {
	return i >= 0 && i <= isa.MaxWord
}
