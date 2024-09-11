package asm

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/runtime/isa"
	str "github.com/kode4food/ale/internal/strings"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/comb/basics"
)

const (
	// ErrUnknownLocalType is raised when a local or private is declared that
	// doesn't have a proper disposition (var, ref, rest)
	ErrUnknownLocalType = "unexpected local type: %s, expected: %s"

	// ErrUnexpectedName is raised when a local name is referenced that hasn't
	// been declared as part of the assembler encoder's scope
	ErrUnexpectedName = "unexpected local name: %s"

	// ErrUnexpectedLabel is raised when a jump or cond-jump instruction refers
	// to a label that hasn't been anchored in the assembler block
	ErrUnexpectedLabel = "unexpected label: %s"

	// ErrExpectedType is raised when a value of a certain type is expected,
	// but not provided at the current position
	ErrExpectedType = "expected %s, got: %s"
)

const (
	Resolve    = data.Local(".resolve")
	Evaluate   = data.Local(".eval")
	Const      = data.Local(".const")
	Local      = data.Local(".local")
	Private    = data.Local(".private")
	PushLocals = data.Local(".push-locals")
	PopLocals  = data.Local(".pop-locals")
)

const (
	pairType = "binding pair"
	seqType  = "sequence"
	symType  = "symbol"
	nameType = "name"
	kwdType  = "keyword"
)

var (
	cellTypes = map[data.Keyword]encoder.CellType{
		data.Keyword("val"):  encoder.ValueCell,
		data.Keyword("ref"):  encoder.ReferenceCell,
		data.Keyword("rest"): encoder.RestCell,
	}

	cellTypeNames = makeCellTypeNames()

	callsOnce sync.Once
	calls     namedAsmParsers
)

// Asm provides indirect access to the Encoder's methods and generators
func Asm(e encoder.Encoder, args ...data.Value) {
	c := getCalls()
	p := makeAsmParser(c)
	emit, err := p.parse(data.Vector(args))
	if err != nil {
		panic(err)
	}
	ae := p.wrapEncoder(e)
	if err := emit(ae); err != nil {
		panic(err)
	}
}

func getCalls() namedAsmParsers {
	callsOnce.Do(func() {
		calls = mergeCalls(
			getInstructionCalls(),
			getEncoderCalls(),
		)
	})
	return calls
}

func mergeCalls(maps ...namedAsmParsers) namedAsmParsers {
	res := namedAsmParsers{}
	for _, m := range maps {
		for k, v := range m {
			if _, ok := res[k]; ok {
				panic(debug.ProgrammerError("duplicate entry: %s", k))
			}
			res[k] = v
		}
	}
	return res
}

func getInstructionCalls() namedAsmParsers {
	res := make(namedAsmParsers, len(isa.Effects))
	for oc, effect := range isa.Effects {
		name := data.Local(str.CamelToSnake(oc.String()))
		res[name] = func(oc isa.Opcode, ao isa.ActOn) asmParse {
			return makeEmitCall(oc, ao)
		}(oc, effect.Operand)
	}
	return res
}

func getEncoderCalls() namedAsmParsers {
	return namedAsmParsers{
		Resolve:    parseArgs(Resolve, 1, resolveCall),
		Evaluate:   parseArgs(Evaluate, 1, evaluateCall),
		Const:      parseArgs(Const, 1, constCall),
		PushLocals: parseArgs(PushLocals, 0, pushLocalsCall),
		PopLocals:  parseArgs(PopLocals, 0, popLocalsCall),
		Local:      parseLocalEncoder(Local, publicNamer),
		Private:    parseLocalEncoder(Private, privateNamer),
		ForEach:    parseForEachCall,
	}
}

func resolveCall(_ *asmParser, args ...data.Value) (asmEmit, error) {
	s, err := assertType[data.Symbol](symType, args[0])
	if err != nil {
		return nil, err
	}

	if l, ok := s.(data.Local); ok {
		return func(e *asmEncoder) error {
			generate.Symbol(e, e.resolvePrivate(l))
			return nil
		}, nil
	}

	return func(e *asmEncoder) error {
		generate.Symbol(e, s)
		return nil
	}, nil
}

func evaluateCall(_ *asmParser, args ...data.Value) (asmEmit, error) {
	return func(e *asmEncoder) error {
		if v, ok := e.resolveEncoderArg(args[0]); ok {
			generate.Value(e, v)
			return nil
		}
		generate.Value(e, args[0])
		return nil
	}, nil
}

func constCall(_ *asmParser, args ...data.Value) (asmEmit, error) {
	return func(e *asmEncoder) error {
		if v, ok := e.resolveEncoderArg(args[0]); ok {
			generate.Literal(e, v)
			return nil
		}
		generate.Literal(e, args[0])
		return nil
	}, nil
}

func pushLocalsCall(*asmParser, ...data.Value) (asmEmit, error) {
	return func(e *asmEncoder) error {
		e.PushLocals()
		return nil
	}, nil
}

func popLocalsCall(*asmParser, ...data.Value) (asmEmit, error) {
	return func(e *asmEncoder) error {
		e.PopLocals()
		return nil
	}, nil
}

func parseLocalEncoder(inst data.Local, toName asmToName) asmParse {
	return parseArgs(inst, 2,
		func(p *asmParser, args ...data.Value) (asmEmit, error) {
			l, err := assertType[data.Local](nameType, args[0])
			if err != nil {
				return nil, err
			}
			k, err := assertType[data.Keyword](kwdType, args[1])
			if err != nil {
				return nil, err
			}
			cellType, err := getCellType(k)
			if err != nil {
				return nil, err
			}

			return func(e *asmEncoder) error {
				name, err := toName(e, l)
				if err != nil {
					return err
				}
				e.AddLocal(name, cellType)
				return nil
			}, nil
		},
	)
}

func publicNamer(e *asmEncoder, l data.Local) (data.Local, error) {
	if v, ok := e.resolveEncoderArg(l); ok {
		if res, ok := v.(data.Local); ok {
			return res, nil
		}
		return "", fmt.Errorf(ErrBadNameResolution, v)
	}
	return l, nil
}

func privateNamer(e *asmEncoder, l data.Local) (data.Local, error) {
	p := gen.Local(l)
	e.private[l] = p
	return p, nil
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

func getCellType(k data.Keyword) (encoder.CellType, error) {
	res, ok := cellTypes[k]
	if ok {
		return res, nil
	}
	return 0, fmt.Errorf(ErrUnknownLocalType, k, cellTypeNames)
}

func makeCellTypeNames() string {
	res := basics.MapKeys(cellTypes)
	var buf strings.Builder
	for i, s := range res {
		switch {
		case i == len(res)-1:
			buf.WriteString(" or ")
		case i != 0:
			buf.WriteString(" ")
		}
		buf.WriteString(s.String())
	}
	return buf.String()
}

func assertType[T data.Value](expected string, val data.Value) (T, error) {
	res, ok := val.(T)
	if !ok {
		var zero T
		return zero, typeError(expected, val)
	}
	return res, nil
}

func typeError(expected string, val data.Value) error {
	return fmt.Errorf(ErrExpectedType, expected, data.ToString(val))
}
