package core

import (
	"errors"
	"fmt"
	"maps"
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
	asmEncoder struct {
		encoder.Encoder
		args    map[data.Local]data.Value
		labels  map[data.Local]isa.Operand
		private map[data.Local]data.Local
	}

	callMap map[data.Local]*call

	call struct {
		special.Call
		argCount int
		hasBlock bool
	}

	toOperandFunc func(data.Value) (isa.Operand, error)
	toNameFunc    func(*asmEncoder, data.Local) data.Local
)

const (
	// ErrUnknownDirective is raised when an unknown directive is called
	ErrUnknownDirective = "unknown directive: %s"

	// ErrUnexpectedForm is raised when an unexpected form is encountered in
	// the assembler block
	ErrUnexpectedForm = "unexpected form: %s"

	// ErrIncompleteInstruction is raised when an instruction is encountered in
	// the assembler block not accompanied by a required operand
	ErrIncompleteInstruction = "incomplete instruction: %s"

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

	// ErrExpectedBinding is raised when a binding vector is expected but not
	// provded to the .for-each call
	ErrExpectedBinding = "expected binding vector, got: %s"

	// ErrExpectedEndOfBlock is raised when an end-of-block marker is expected
	// but an end of stream is encountered instead
	ErrExpectedEndOfBlock = "expected end of block"
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

var (
	gen = data.NewSymbolGenerator()

	cellTypes = map[data.Keyword]encoder.CellType{
		data.Keyword("val"):  encoder.ValueCell,
		data.Keyword("ref"):  encoder.ReferenceCell,
		data.Keyword("rest"): encoder.RestCell,
	}

	cellTypeNames = makeCellTypeNames()

	callsOnce sync.Once
	calls     callMap
)

// Asm provides indirect access to the Encoder's methods and generators
func Asm(e encoder.Encoder, args ...data.Value) {
	makeAsmEncoder(e).process(data.Vector(args))
}

func makeAsmEncoder(e encoder.Encoder) *asmEncoder {
	return &asmEncoder{
		Encoder: e,
		labels:  map[data.Local]isa.Operand{},
		args:    map[data.Local]data.Value{},
		private: map[data.Local]data.Local{},
	}
}

func (e *asmEncoder) withParams(n data.Locals, v data.Vector) *asmEncoder {
	args := make(map[data.Local]data.Value, len(n))
	for i, k := range n {
		args[k] = v[i]
	}
	res := e.copy()
	res.args = args
	return res
}

func (e *asmEncoder) copy() *asmEncoder {
	res := *e
	res.args = maps.Clone(res.args)
	res.private = maps.Clone(res.private)
	res.labels = maps.Clone(res.labels)
	return &res
}

func (e *asmEncoder) process(forms data.Sequence) {
	if f, r, ok := forms.Split(); ok {
		if l, ok := f.(data.Local); ok {
			switch l {
			case MakeSpecial:
				e.makeSpecialCall(r)
				return
			}
		}
	}
	e.encode(forms)
}

func (e *asmEncoder) makeSpecialCall(forms data.Sequence) {
	pc := parseParamCases(forms)
	e.makeSpecialFromCases(pc)
}

func (e *asmEncoder) makeSpecialFromCases(pc *paramCases) {
	ac := pc.makeChecker()
	f := pc.makeFetchers()
	fn := func(e encoder.Encoder, args ...data.Value) {
		if err := ac(len(args)); err != nil {
			panic(err)
		}
		for i, c := range pc.Cases() {
			if a, ok := f[i](args); ok {
				ae := makeAsmEncoder(e).withParams(c.params, a)
				ae.encode(c.body)
				return
			}
		}
		panic(errors.New(ErrNoMatchingParamPattern))
	}
	e.Emit(isa.Const, e.AddConstant(special.Call(fn)))
}

func (e *asmEncoder) encode(forms data.Sequence) {
	for f, r, ok := forms.Split(); ok; f, r, ok = r.Split() {
		switch name := f.(type) {
		case data.Keyword:
			e.Emit(isa.Label, e.getLabelIndex(name.Name()))
		case data.Local:
			d, ok := getCalls()[name]
			if !ok {
				panic(fmt.Errorf(ErrUnknownDirective, name))
			}
			args, rest, ok := take(r, d.argCount)
			if !ok {
				panic(fmt.Errorf(ErrIncompleteInstruction, name))
			}
			if d.hasBlock {
				var block data.Vector
				block, rest = parseBlock(rest)
				args = append(args, block...)
			}
			d.Call(e, args...)
			r = rest
		default:
			panic(fmt.Errorf(ErrUnexpectedForm, data.ToString(f)))
		}
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

func (e *asmEncoder) toOperands(oc isa.Opcode, args data.Vector) []isa.Operand {
	return basics.Map(args, func(a data.Value) isa.Operand {
		ao := isa.Effects[oc].Operand
		toOperand := e.getToOperandFor(ao)
		r, err := toOperand(a)
		if err != nil {
			panic(err)
		}
		return r
	})
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
		if v, ok := e.resolveEncoderArg(val); ok {
			val = v
		}
		if val, ok := val.(data.Keyword); ok {
			return e.getLabelIndex(val.Name()), nil
		}
		return toOperand(val)
	}, ErrUnexpectedLabel)
}

func (e *asmEncoder) makeNameToWord() toOperandFunc {
	return wrapToOperandError(func(val data.Value) (isa.Operand, error) {
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
		return toOperand(val)
	}, ErrUnexpectedName)
}

func (e *asmEncoder) resolveEncoderArg(v data.Value) (data.Value, bool) {
	if v, ok := v.(data.Local); ok {
		res, ok := e.args[v]
		return res, ok
	}
	return nil, false
}

func (e *asmEncoder) resolvePrivate(l data.Local) data.Local {
	if g, ok := e.private[l]; ok {
		return g
	}
	return l
}

func getCalls() callMap {
	callsOnce.Do(func() {
		calls = mergeCalls(
			getInstructionCalls(),
			getEncoderCalls(),
		)
	})
	return calls
}

func mergeCalls(maps ...callMap) callMap {
	res := callMap{}
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

func getInstructionCalls() callMap {
	res := make(callMap, len(isa.Effects))
	for oc, effect := range isa.Effects {
		name := data.Local(str.CamelToSnake(oc.String()))
		res[name] = func(oc isa.Opcode, ao isa.ActOn) *call {
			return makeEmitCall(oc, ao)
		}(oc, effect.Operand)
	}
	return res
}

func getEncoderCalls() callMap {
	return callMap{
		Resolve:    {Call: resolveCall, argCount: 1},
		Evaluate:   {Call: evaluateCall, argCount: 1},
		ForEach:    {Call: forEachCall, argCount: 1, hasBlock: true},
		Const:      {Call: constCall, argCount: 1},
		PushLocals: {Call: pushLocalsCall},
		PopLocals:  {Call: popLocalsCall},
		Local:      {Call: makeLocalEncoder(publicNamer), argCount: 2},
		Private:    {Call: makeLocalEncoder(privateNamer), argCount: 2},
	}
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

func resolveCall(e encoder.Encoder, args ...data.Value) {
	s := args[0].(data.Symbol)
	if l, ok := s.(data.Local); ok {
		generate.Symbol(e, e.(*asmEncoder).resolvePrivate(l))
		return
	}
	generate.Symbol(e, s)
}

func evaluateCall(e encoder.Encoder, args ...data.Value) {
	if v, ok := e.(*asmEncoder).resolveEncoderArg(args[0]); ok {
		generate.Value(e, v)
		return
	}
	generate.Value(e, args[0])
}

func forEachCall(e encoder.Encoder, args ...data.Value) {
	n, v, ok := parseForEachBinding(args[0])
	if !ok {
		panic(fmt.Errorf(ErrExpectedBinding, args[0]))
	}
	s, ok := e.(*asmEncoder).resolveEncoderArg(v)
	if !ok {
		panic(fmt.Errorf(ErrUnexpectedParameter, v))
	}
	ae := e.(*asmEncoder).copy()
	forms := data.Vector(args[1:])
	seq := s.(data.Sequence)
	for f, r, ok := seq.Split(); ok; f, r, ok = r.Split() {
		ae.args[n] = f
		ae.encode(forms)
	}
}

func parseForEachBinding(v data.Value) (data.Local, data.Value, bool) {
	b, ok := v.(data.Vector)
	if !ok || len(b) != 2 {
		return "", nil, false
	}
	l, ok := b[0].(data.Local)
	if !ok {
		return "", nil, false
	}
	return l, b[1], true
}

func parseBlock(s data.Sequence) (data.Vector, data.Sequence) {
	var res data.Vector
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		if l, ok := f.(data.Local); ok && l == EndBlock {
			return res, r
		}
		res = append(res, f)
	}
	panic(errors.New(ErrExpectedEndOfBlock))
}

func constCall(e encoder.Encoder, args ...data.Value) {
	if v, ok := e.(*asmEncoder).resolveEncoderArg(args[0]); ok {
		generate.Literal(e, v)
		return
	}
	generate.Literal(e, args[0])
}

func pushLocalsCall(e encoder.Encoder, _ ...data.Value) {
	e.PushLocals()
}

func popLocalsCall(e encoder.Encoder, _ ...data.Value) {
	e.PopLocals()
}

func makeLocalEncoder(toName toNameFunc) special.Call {
	return func(e encoder.Encoder, args ...data.Value) {
		name := toName(e.(*asmEncoder), args[0].(data.Local))
		kwd := args[1].(data.Keyword)
		cellType, ok := cellTypes[kwd]
		if !ok {
			panic(fmt.Errorf(ErrUnknownLocalType, kwd, cellTypeNames))
		}
		e.AddLocal(name, cellType)
	}
}

func publicNamer(e *asmEncoder, l data.Local) data.Local {
	if v, ok := e.resolveEncoderArg(l); ok {
		if res, ok := v.(data.Local); ok {
			return res
		}
		panic(fmt.Errorf(ErrBadNameResolution, v))
	}
	return l
}

func privateNamer(e *asmEncoder, l data.Local) data.Local {
	p := gen.Local(l)
	e.private[l] = p
	return p
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
