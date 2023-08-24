package special

import (
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/strings"
	"github.com/kode4food/ale/runtime/isa"
)

type (
	call struct {
		encoder.Call
		argCount int
	}

	callMap map[data.Name]*call
)

// Error messages
const (
	ErrIncompleteInstruction = "incomplete instruction: %s"
	ErrExpectedWord          = "expected unsigned word: %s"
	ErrUnknownDirective      = "unknown opcode: %s"
	ErrUnexpectedForm        = "unexpected form: %s"
)

// Asm provides indirect access to the Encoder's methods and generators
//
// (asm*
//   push-locals
//   true
//   cond-jump some-label
//   ret-nil
// :some-label
//   literal "hello"
//   pop-locals)

var (
	instCalls = getInstructionCalls()
	encCalls  = getEncoderCalls()
	calls     = mergeCalls(instCalls, encCalls)
)

func Asm(e encoder.Encoder, args ...data.Value) {
	labels := map[data.Name]isa.Index{}
	v := data.NewVector(args...)
	for f, r, ok := v.Split(); ok; f, r, ok = r.Split() {
		switch v := f.(type) {
		case data.Keyword:
			n := v.Name()
			if _, ok := labels[n]; ok {
				panic(fmt.Errorf("label already placed: %s", n))
			}
			labels[n] = e.NewLabel()
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
			e.Emit(oc, toWords(args)...)
		},
		argCount: argCount,
	}
}

func getEncoderCalls() callMap {
	return callMap{
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

func toWords(args data.Values) []isa.Coder {
	words := make([]isa.Coder, len(args))
	for i, a := range args {
		arg := a.(data.Integer)
		if !isValidWord(arg) {
			panic(fmt.Errorf(ErrExpectedWord, args[i]))
		}
		words = append(words, isa.Word(arg))
	}
	return words
}

func isValidWord(i data.Integer) bool {
	return i >= 0 && i <= isa.MaxWord
}
