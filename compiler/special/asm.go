package special

import (
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
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
)

// Error messages
const (
	ErrUnknownDirective         = "unknown directive: %s"
	ErrUnexpectedForm           = "unexpected form: %s"
	ErrIncompleteInstruction    = "incomplete instruction: %s"
	ErrUnexpectedInstructionArg = "unexpected instruction argument: %s"
	ErrUnexpectedName           = "unexpected local name: %s"
	ErrUnexpectedLabel          = "unexpected label: %s"
	ErrExpectedWord             = "expected unsigned word: %s"
)

var (
	instCalls = getInstructionCalls()
	encCalls  = getEncoderCalls()
	genCalls  = getGeneratorCalls()
	calls     = mergeCalls(instCalls, encCalls, genCalls)

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
		".const": {
			Call: func(e encoder.Encoder, args ...data.Value) {
				index := e.AddConstant(args[0])
				e.Emit(isa.Const, index)
			},
			argCount: 1,
		},
		".local": {
			Call: func(e encoder.Encoder, args ...data.Value) {
				name := args[0].(data.LocalSymbol).Name()
				kwd := args[1].(data.Keyword)
				cellType, ok := cellTypes[kwd]
				if !ok {
					panic(fmt.Errorf("unknown local type: %s", kwd))
				}
				e.AddLocal(name, cellType)
			},
			argCount: 2,
		},
		".push-locals": {
			Call: func(e encoder.Encoder, _ ...data.Value) {
				e.PushLocals()
			},
		},
		".pop-locals": {
			Call: func(e encoder.Encoder, _ ...data.Value) {
				e.PopLocals()
			},
		},
	}
}

func getGeneratorCalls() callMap {
	return callMap{}
}

func mergeCalls(maps ...callMap) callMap {
	res := callMap{}
	for _, m := range maps {
		for k, v := range m {
			if _, ok := res[k]; ok {
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
	words := make([]isa.Coder, len(args))
	for i, a := range args {
		switch arg := a.(type) {
		case data.Integer:
			if !isValidWord(arg) {
				panic(fmt.Errorf(ErrExpectedWord, a))
			}
			words[i] = isa.Word(arg)
		case data.Keyword:
			if !isa.Effects[oc].Labels {
				panic(fmt.Errorf(ErrUnexpectedLabel, a))
			}
			words[i] = e.getLabelIndex(arg)
		case data.LocalSymbol:
			cell, ok := e.ResolveLocal(arg.Name())
			if !ok || !isa.Effects[oc].Locals {
				panic(fmt.Errorf(ErrUnexpectedName, arg))
			}
			words[i] = cell.Index
		default:
			panic(fmt.Errorf(ErrUnexpectedInstructionArg, a.String()))
		}
	}
	return words
}

func isValidWord(i data.Integer) bool {
	return i >= 0 && i <= isa.MaxWord
}
