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
		labels map[data.Name]isa.Operand
		args   map[data.Name]data.Value
	}

	call struct {
		encoder.Call
		argCount int
	}

	callMap map[data.Name]*call

	toOperandFunc func(data.Value) (isa.Operand, error)
)

// Error messages
const (
	ErrUnknownDirective      = "unknown directive: %s"
	ErrUnexpectedForm        = "unexpected form: %s"
	ErrIncompleteInstruction = "incomplete instruction: %s"
	ErrUnknownLocalType      = "unknown local type: %s"
	ErrUnexpectedName        = "unexpected local name: %s"
	ErrUnexpectedLabel       = "unexpected label: %s"
	ErrExpectedOperand       = "expected unsigned word: %s"
)

const (
	MakeEncoder = data.Name("!make-encoder")
	Value       = data.Name(".value")
	Const       = data.Name(".const")
	Local       = data.Name(".local")
	PushLocals  = data.Name(".push-locals")
	PopLocals   = data.Name(".pop-locals")
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
	makeAsmEncoder(e).process(data.NewVector(args...))
}

func makeAsmEncoder(e encoder.Encoder) *asmEncoder {
	return &asmEncoder{
		Encoder: e,
		labels:  map[data.Name]isa.Operand{},
		args:    map[data.Name]data.Value{},
	}
}

func (e *asmEncoder) withParams(n data.Names, v data.Values) *asmEncoder {
	args := make(map[data.Name]data.Value, len(n))
	for i, k := range n {
		args[k] = v[i]
	}
	res := *e
	res.args = args
	return &res
}

func (e *asmEncoder) process(forms data.Sequence) {
	if f, r, ok := forms.Split(); ok {
		if l, ok := f.(data.LocalSymbol); ok {
			switch l.Name() {
			case MakeEncoder:
				e.makeEncoderCall(r)
				return
			}
		}
	}
	e.encode(forms)
}

func (e *asmEncoder) makeEncoderCall(forms data.Sequence) {
	cases := parseParamCases(forms)
	ac := cases.makeArityChecker()
	f := cases.makeFetchers()
	fn := func(e encoder.Encoder, args ...data.Value) {
		if err := ac(len(args)); err != nil {
			panic(err)
		}
		for i, c := range cases {
			if a, ok := f[i](args); ok {
				ae := makeAsmEncoder(e).withParams(c.params, a)
				ae.encode(c.body)
				return
			}
		}
		panic(ErrNoMatchingParamPattern)
	}
	e.Emit(isa.Const, e.AddConstant(encoder.Call(fn)))
}

func (e *asmEncoder) encode(forms data.Sequence) {
	for f, r, ok := forms.Split(); ok; f, r, ok = r.Split() {
		switch v := f.(type) {
		case data.Keyword:
			e.Emit(isa.Label, e.getLabelIndex(v))
		case data.LocalSymbol:
			n := v.Name()
			if d, ok := calls[n]; ok {
				if args, rest, ok := take(r, d.argCount); ok {
					d.Call(e, args...)
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

func (e *asmEncoder) getLabelIndex(k data.Keyword) isa.Operand {
	name := k.Name()
	if idx, ok := e.labels[name]; ok {
		return idx
	}
	idx := e.NewLabel()
	e.labels[name] = idx
	return idx
}

func (e *asmEncoder) toOperands(
	oc isa.Opcode, args data.Values,
) []isa.Operand {
	res := make([]isa.Operand, len(args))
	for i, a := range args {
		ao := isa.Effects[oc].Operand
		toOperand := e.getToOperandFor(ao)
		r, err := toOperand(a)
		if err != nil {
			panic(err)
		}
		res[i] = r
	}
	return res
}

func (e *asmEncoder) getToOperandFor(ao isa.ActOn) toOperandFunc {
	switch ao {
	case isa.Locals:
		return e.makeNameToWord()
	case isa.Labels:
		return e.makeLabelToWord()
	default:
		return toOperand
	}
}

func (e *asmEncoder) makeLabelToWord() toOperandFunc {
	return wrapToOperandError(func(val data.Value) (isa.Operand, error) {
		if val, ok := val.(data.Keyword); ok {
			return e.getLabelIndex(val), nil
		}
		return toOperand(val)
	}, ErrUnexpectedLabel)
}

func (e *asmEncoder) makeNameToWord() toOperandFunc {
	return wrapToOperandError(func(val data.Value) (isa.Operand, error) {
		if val, ok := val.(data.LocalSymbol); ok {
			if cell, ok := e.ResolveLocal(val.Name()); ok {
				return cell.Index, nil
			}
			return 0, fmt.Errorf(ErrUnexpectedName, val)
		}
		return toOperand(val)
	}, ErrUnexpectedName)
}

func getInstructionCalls() callMap {
	res := make(callMap, len(isa.Effects))
	for oc, effect := range isa.Effects {
		name := data.Name(strings.CamelToSnake(oc.String()))
		res[name] = func(oc isa.Opcode, ao isa.ActOn) *call {
			return makeEmitCall(oc, ao)
		}(oc, effect.Operand)
	}
	return res
}

func makeEmitCall(oc isa.Opcode, actOn isa.ActOn) *call {
	argCount := 0
	if actOn != isa.Nothing {
		argCount = 1
	}
	return &call{
		Call: func(e encoder.Encoder, args ...data.Value) {
			e.Emit(oc, e.(*asmEncoder).toOperands(oc, args)...)
		},
		argCount: argCount,
	}
}

func getEncoderCalls() callMap {
	return callMap{
		Value: {
			Call: func(e encoder.Encoder, args ...data.Value) {
				if arg, ok := args[0].(data.LocalSymbol); ok {
					if v, ok := e.(*asmEncoder).args[arg.Name()]; ok {
						generate.Value(e, v)
						return
					}
				}
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

func wrapToOperandError(toOperand toOperandFunc, errStr string) toOperandFunc {
	return func(val data.Value) (isa.Operand, error) {
		res, err := toOperand(val)
		if err != nil {
			return 0, errors.Join(fmt.Errorf(errStr, val), err)
		}
		return res, nil
	}
}

func toOperand(val data.Value) (isa.Operand, error) {
	if val, ok := val.(data.Integer); ok {
		if isValidOperand(val) {
			return isa.Operand(val), nil
		}
	}
	return 0, fmt.Errorf(ErrExpectedOperand, val)
}

func isValidOperand(i data.Integer) bool {
	return i >= 0 && i <= isa.OperandMask
}
