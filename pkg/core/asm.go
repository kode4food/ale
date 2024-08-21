package core

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/compiler/special"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/runtime/isa"
	str "github.com/kode4food/ale/internal/strings"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/comb/basics"
)

type (
	asmParser struct {
		params  data.Locals
		private map[data.Local]data.Local
	}

	asmEncoder struct {
		encoder.Encoder
		*asmParser
		args   map[data.Local]data.Value
		labels map[data.Local]isa.Operand
	}

	namedAsmParsers map[data.Local]asmParse

	asmParse     func(*asmParser, data.Sequence) (asmEmit, data.Sequence, error)
	asmArgsParse func(*asmParser, ...data.Value) (asmEmit, error)

	asmEmit      func(*asmEncoder) error
	asmToOperand func(*asmEncoder, data.Value) (isa.Operand, error)
	asmToName    func(*asmEncoder, data.Local) (data.Local, error)
)

const (
	// ErrUnknownDirective is raised when an unknown directive is called
	ErrUnknownDirective = "unknown directive: %s"

	// ErrUnexpectedForm is raised when an unexpected form is encountered in
	// the assembler block
	ErrUnexpectedForm = "unexpected form: %s"

	// ErrIncompleteInstruction is raised when an instruction is encountered in
	// the assembler block not accompanied by enough operands
	ErrIncompleteInstruction = "incomplete %s instruction, args expected: %d"

	// ErrUnknownLocalType is raised when a local or private is declared that
	// doesn't have a proper disposition (var, ref, rest)
	ErrUnknownLocalType = "unexpected local type: %s, expected: %s"

	// ErrUnexpectedParameter is raised when an encoder parameter is not found.
	// These are declared using the !make-special directive
	ErrUnexpectedParameter = "unexpected parameter name: %s"

	// ErrUnexpectedName is raised when a local name is referenced that hasn't
	// been declared as part of the assembler encoder's scope
	ErrUnexpectedName = "unexpected local name: %s"

	// ErrUnexpectedLabel is raised when a jump or cond-jump instruction refers
	// to a label that hasn't been anchored in the assembler block
	ErrUnexpectedLabel = "unexpected label: %s"

	// ErrBadNameResolution is raised when an attempt is made to bind a local
	// using an argument to the encoder that is not a Local symbol
	ErrBadNameResolution = "encoder argument is not a name: %s"

	// ErrExpectedEndOfBlock is raised when an end-of-block marker is expected
	// but an end of stream is encountered instead
	ErrExpectedEndOfBlock = "expected end of block"

	// ErrExpectedType is raised when a value of a certain type is expected,
	// but not provided at the current position
	ErrExpectedType = "expected %s, got: %s"

	// ErrPairExpected is raised when a vector is provided, but does not
	// contain exactly two elements
	ErrPairExpected = "binding pair expected, got %d elements"
)

const (
	MakeSpecial = data.Local("!make-special")
	Resolve     = data.Local(".resolve")
	Evaluate    = data.Local(".eval")
	ForEach     = data.Local(".for-each")
	Const       = data.Local(".const")
	Local       = data.Local(".local")
	Private     = data.Local(".private")
	PushLocals  = data.Local(".push-locals")
	PopLocals   = data.Local(".pop-locals")
	EndBlock    = data.Local(".end")
)

const (
	pairType = "binding pair"
	seqType  = "sequence"
	symType  = "symbol"
	nameType = "name"
	kwdType  = "keyword"
)

var (
	gen = data.NewSymbolGenerator()

	cellTypes = map[data.Keyword]encoder.CellType{
		data.Keyword("val"):  encoder.ValueCell,
		data.Keyword("ref"):  encoder.ReferenceCell,
		data.Keyword("rest"): encoder.RestCell,
	}

	cellTypeNames = makeCellTypeNames()

	callsOnce sync.Once
	calls     namedAsmParsers
)

func noAsmEmit(*asmEncoder) error { return nil }

// Asm provides indirect access to the Encoder's methods and generators
func Asm(e encoder.Encoder, args ...data.Value) {
	p := makeAsmParser()
	emit, err := p.parse(data.Vector(args))
	if err != nil {
		panic(err)
	}
	ae := p.wrapEncoder(e)
	if err := emit(ae); err != nil {
		panic(err)
	}
}

func makeAsmParser() *asmParser {
	return &asmParser{
		params:  data.Locals{},
		private: map[data.Local]data.Local{},
	}
}

func (p *asmParser) withParams(n data.Locals) *asmParser {
	res := p.copy()
	res.params = slices.Clone(n)
	return res
}

func (p *asmParser) copy() *asmParser {
	res := *p
	res.params = slices.Clone(res.params)
	res.private = maps.Clone(res.private)
	return &res
}

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
	}
}

func (p *asmParser) parse(forms data.Sequence) (asmEmit, error) {
	if f, r, ok := forms.Split(); ok && f == MakeSpecial {
		return p.specialCall(r)
	}
	return p.sequence(forms)
}

func (p *asmParser) specialCall(forms data.Sequence) (asmEmit, error) {
	pc := parseParamCases(forms)
	cases := pc.Cases()
	ap := make([]*asmParser, len(cases))
	emitters := make([]asmEmit, len(cases))
	for i, c := range cases {
		ap[i] = p.withParams(c.params)
		e, err := ap[i].sequence(c.body)
		if err != nil {
			return nil, err
		}
		emitters[i] = e
	}

	ac := pc.makeChecker()
	fetchers := pc.makeFetchers()

	fn := func(e encoder.Encoder, args ...data.Value) {
		if err := ac(len(args)); err != nil {
			panic(err)
		}
		for i, f := range fetchers {
			if a, ok := f(args); ok {
				ae := ap[i].wrapEncoder(e, a...)
				if err := emitters[i](ae); err != nil {
					panic(err)
				}
				return
			}
		}
		panic(errors.New(ErrNoMatchingParamPattern))
	}

	return func(e *asmEncoder) error {
		e.Emit(isa.Const, e.AddConstant(special.Call(fn)))
		return nil
	}, nil
}

func (p *asmParser) next(s data.Sequence) (asmEmit, data.Sequence, error) {
	f, r, _ := s.Split()
	switch t := f.(type) {
	case data.Keyword:
		return func(e *asmEncoder) error {
			e.Emit(isa.Label, e.getLabelIndex(t.Name()))
			return nil
		}, r, nil
	case data.Local:
		parse, ok := getCalls()[t]
		if !ok {
			return nil, nil, fmt.Errorf(ErrUnknownDirective, t)
		}
		return parse(p, r)
	default:
		return nil, nil, fmt.Errorf(ErrUnexpectedForm, data.ToString(f))
	}
}

func (p *asmParser) sequence(s data.Sequence) (asmEmit, error) {
	if s.IsEmpty() {
		return noAsmEmit, nil
	}
	next, rest, err := p.next(s)
	if err != nil {
		return nil, err
	}
	return p.rest(next, rest)
}

func (p *asmParser) rest(emit asmEmit, r data.Sequence) (asmEmit, error) {
	next, err := p.sequence(r)
	if err != nil {
		return nil, err
	}
	return func(e *asmEncoder) error {
		if err := emit(e); err != nil {
			return err
		}
		return next(e)
	}, nil
}

func (p *asmParser) block(s data.Sequence) (asmEmit, data.Sequence, error) {
	if s.IsEmpty() {
		return nil, nil, errors.New(ErrExpectedEndOfBlock)
	}
	f, r, _ := s.Split()
	if f == EndBlock {
		return noAsmEmit, r, nil
	}
	next, rest, err := p.next(s)
	if err != nil {
		return nil, nil, err
	}
	return p.blockRest(next, rest)
}

func (p *asmParser) blockRest(
	emit asmEmit, r data.Sequence,
) (asmEmit, data.Sequence, error) {
	next, rest, err := p.block(r)
	if err != nil {
		return nil, nil, err
	}
	return func(e *asmEncoder) error {
		if err := emit(e); err != nil {
			return err
		}
		return next(e)
	}, rest, nil
}

func (p *asmParser) getToOperandFor(ao isa.ActOn) asmToOperand {
	switch ao {
	case isa.Locals:
		return p.makeNameToWord()
	case isa.Labels:
		return p.makeLabelToWord()
	default:
		return toOperand
	}
}

func (p *asmParser) makeLabelToWord() asmToOperand {
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

func (p *asmParser) makeNameToWord() asmToOperand {
	return wrapToOperandError(ErrUnexpectedName,
		func(e *asmEncoder, val data.Value) (isa.Operand, error) {
			if v, ok := e.resolveEncoderArg(val); ok {
				val = v
			}
			if val, ok := val.(data.Local); ok {
				n := p.resolvePrivate(val)
				if cell, ok := e.ResolveLocal(n); ok {
					return cell.Index, nil
				}
				return 0, fmt.Errorf(ErrUnexpectedName, val)
			}
			return toOperand(e, val)
		},
	)
}

func (p *asmParser) resolvePrivate(l data.Local) data.Local {
	if g, ok := p.private[l]; ok {
		return g
	}
	return l
}

func (e *asmEncoder) getLabelIndex(n data.Local) isa.Operand {
	if idx, ok := e.labels[n]; ok {
		return idx
	}
	idx := e.NewLabel()
	e.labels[n] = idx
	return idx
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
		ForEach:    parseForEachCall,
		Const:      parseArgs(Const, 1, constCall),
		PushLocals: parseArgs(PushLocals, 0, pushLocalsCall),
		PopLocals:  parseArgs(PopLocals, 0, popLocalsCall),
		Local:      parseLocalEncoder(Local, publicNamer),
		Private:    parseLocalEncoder(Private, privateNamer),
	}
}

func parseArgs(inst data.Local, argsLen int, fn asmArgsParse) asmParse {
	return func(p *asmParser, s data.Sequence) (asmEmit, data.Sequence, error) {
		args, rest, ok := take(s, argsLen)
		if !ok {
			return nil, nil, fmt.Errorf(ErrIncompleteInstruction, inst, argsLen)
		}
		res, err := fn(p, args...)
		if err != nil {
			return nil, nil, err
		}
		return res, rest, nil
	}
}

func resolveCall(_ *asmParser, args ...data.Value) (asmEmit, error) {
	s, err := assertType[data.Symbol](symType, args[0])
	if err != nil {
		return nil, err
	}

	return func(e *asmEncoder) error {
		if l, ok := s.(data.Local); ok {
			generate.Symbol(e, e.resolvePrivate(l))
			return nil
		}
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

func parseForEachCall(
	p *asmParser, s data.Sequence,
) (asmEmit, data.Sequence, error) {
	f, r, ok := s.Split()
	if !ok {
		return nil, nil, typeError(pairType, f)
	}
	k, v, err := parseForEachBinding(f)
	if err != nil {
		return nil, nil, err
	}
	pc := p.withParams(data.Locals{k})
	block, rest, err := p.block(r)
	if err != nil {
		return nil, nil, err
	}

	return func(e *asmEncoder) error {
		s, ok := e.resolveEncoderArg(v)
		if !ok {
			return fmt.Errorf(ErrUnexpectedParameter, v)
		}
		seq, err := assertType[data.Sequence](seqType, s)
		if err != nil {
			return err
		}
		for f, r, ok := seq.Split(); ok; f, r, ok = r.Split() {
			if err := block(pc.wrapEncoder(e, f)); err != nil {
				return err
			}
		}
		return nil
	}, rest, nil
}

func parseForEachBinding(v data.Value) (data.Local, data.Value, error) {
	b, err := assertType[data.Vector](pairType, v)
	if err != nil {
		return "", nil, err
	}
	if len(b) != 2 {
		return "", nil, fmt.Errorf(ErrPairExpected, len(b))
	}
	l, err := assertType[data.Local](nameType, b[0])
	if err != nil {
		return "", nil, err
	}
	return l, b[1], nil
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
			cellType, ok := cellTypes[k]
			if !ok {
				return nil, fmt.Errorf(ErrUnknownLocalType, k, cellTypeNames)
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

func take(s data.Sequence, count int) (data.Vector, data.Sequence, bool) {
	var f data.Value
	var ok bool
	res := make(data.Vector, count)
	for i := 0; i < count; i++ {
		if f, s, ok = s.Split(); !ok {
			return nil, nil, false
		}
		res[i] = f
	}
	return res, s, true
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
