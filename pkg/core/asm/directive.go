package asm

import (
	"fmt"
	"strings"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/comb/basics"
)

const (
	// ErrBadNameResolution is raised when an attempt is made to bind a local
	// using an argument to the encoder that is not a Local symbol
	ErrBadNameResolution = "encoder argument is not a name: %s"

	// ErrUnknownLocalType is raised when a local or private is declared that
	// doesn't have a proper disposition (var, ref, rest)
	ErrUnknownLocalType = "unexpected local type: %s, expected: %s"

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
)

func getDirectiveCalls() namedAsmParsers {
	return namedAsmParsers{
		Const:      parseArgs(Const, 1, constCall),
		Evaluate:   parseArgs(Evaluate, 1, evaluateCall),
		ForEach:    parseForEachCall,
		Local:      parseLocalEncoder(Local, publicNamer),
		PopLocals:  parseArgs(PopLocals, 0, popLocalsCall),
		Private:    parseLocalEncoder(Private, privateNamer),
		PushLocals: parseArgs(PushLocals, 0, pushLocalsCall),
		Resolve:    parseArgs(Resolve, 1, resolveCall),
	}
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

func publicNamer(e *asmEncoder, l data.Local) (data.Local, error) {
	if v, ok := e.resolveEncoderArg(l); ok {
		if res, ok := v.(data.Local); ok {
			return res, nil
		}
		return "", fmt.Errorf(ErrBadNameResolution, v)
	}
	return l, nil
}

func popLocalsCall(*asmParser, ...data.Value) (asmEmit, error) {
	return func(e *asmEncoder) error {
		e.PopLocals()
		return nil
	}, nil
}

func privateNamer(e *asmEncoder, l data.Local) (data.Local, error) {
	p := gen.Local(l)
	e.private[l] = p
	return p, nil
}

func pushLocalsCall(*asmParser, ...data.Value) (asmEmit, error) {
	return func(e *asmEncoder) error {
		e.PushLocals()
		return nil
	}, nil
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