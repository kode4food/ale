package builtin

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/internal/strings"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/comb/basics"
)

type (
	asmEncoder struct {
		encoder.Encoder
		labels  map[data.Local]isa.Operand
		args    map[data.Local]data.Value
		private map[data.Local]data.Local
	}

	call struct {
		special.Call
		argCount int
	}

	callMap map[data.Local]*call

	toOperandFunc func(data.Value) (isa.Operand, error)
)

const (
	// ErrUnknownDirective is raised when an unknown directive is called
	ErrUnknownDirective = "unknown directive: %s"

	// ErrUnexpectedForm is raised when an unexpected form is encountered in
	// the assembler block
	ErrUnexpectedForm = "unexpected form: %s"

	// ErrIncompleteInstruction is raised when an instruction is encountered in
	// the assembler block that isn't accompanied by a required operand
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
)

var (
	instructionCalls = getInstructionCalls()
	encoderCalls     = getEncoderCalls()
	calls            = mergeCalls(instructionCalls, encoderCalls)

	gen = data.NewSymbolGenerator()

	cellTypes = map[data.Keyword]encoder.CellType{
		data.Keyword("val"):  encoder.ValueCell,
		data.Keyword("ref"):  encoder.ReferenceCell,
		data.Keyword("rest"): encoder.RestCell,
	}

	cellTypeNames = makeCellTypeNames()
)

// Asm provides indirect access to the Encoder's methods and generators
func Asm(e encoder.Encoder, args ...data.Value) {
	makeAsmEncoder(e).process(data.NewVector(args...))
}

func makeAsmEncoder(e encoder.Encoder) *asmEncoder {
	return &asmEncoder{
		Encoder: e,
		labels:  map[data.Local]isa.Operand{},
		args:    map[data.Local]data.Value{},
		private: map[data.Local]data.Local{},
	}
}

func (e *asmEncoder) withParams(n data.Locals, v data.Values) *asmEncoder {
	args := make(map[data.Local]data.Value, len(n))
	for i, k := range n {
		args[k] = v[i]
	}
	res := *e
	res.args = args
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
	e.Emit(isa.Const, e.AddConstant(special.Call(fn)))
}

func (e *asmEncoder) encode(forms data.Sequence) {
	for f, r, ok := forms.Split(); ok; f, r, ok = r.Split() {
		switch name := f.(type) {
		case data.Keyword:
			e.Emit(isa.Label, e.getLabelIndex(name.Name()))
		case data.Local:
			d, ok := calls[name]
			if !ok {
				panic(fmt.Errorf(ErrUnknownDirective, name))
			}
			args, rest, ok := take(r, d.argCount)
			if !ok {
				panic(fmt.Errorf(ErrIncompleteInstruction, name))
			}
			d.Call(e, args...)
			r = rest
		default:
			panic(fmt.Errorf(ErrUnexpectedForm, f.String()))
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

func (e *asmEncoder) toOperands(
	oc isa.Opcode, args data.Values,
) []isa.Operand {
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

func getInstructionCalls() callMap {
	res := make(callMap, len(isa.Effects))
	for oc, effect := range isa.Effects {
		name := data.Local(strings.CamelToSnake(oc.String()))
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
		Resolve:    {Call: resolveCall, argCount: 1},
		Evaluate:   {Call: evaluateCall, argCount: 1},
		ForEach:    {Call: forEachCall, argCount: 2},
		Const:      {Call: constCall, argCount: 1},
		PushLocals: {Call: pushLocalsCall},
		PopLocals:  {Call: popLocalsCall},
		Local:      {Call: makeLocalEncoder(publicNamer), argCount: 2},
		Private:    {Call: makeLocalEncoder(privateNamer), argCount: 2},
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
	name := args[0].(data.Local)
	encode := eval.Value(e.Globals(), args[1]).(special.Call)
	s, ok := e.(*asmEncoder).resolveEncoderArg(name)
	if !ok {
		panic(fmt.Errorf(ErrUnexpectedParameter, name))
	}
	seq := s.(data.Sequence)
	for f, r, ok := seq.Split(); ok; f, r, ok = r.Split() {
		encode(e, f)
	}
}

func constCall(e encoder.Encoder, args ...data.Value) {
	if v, ok := e.(*asmEncoder).resolveEncoderArg(args[0]); ok {
		generate.Literal(e, v)
		return
	}
	index := e.AddConstant(args[0])
	e.Emit(isa.Const, index)
}

func pushLocalsCall(e encoder.Encoder, _ ...data.Value) {
	e.PushLocals()
}

func popLocalsCall(e encoder.Encoder, _ ...data.Value) {
	e.PopLocals()
}

func makeLocalEncoder(
	namer func(e *asmEncoder, l data.Local) data.Local,
) special.Call {
	return func(e encoder.Encoder, args ...data.Value) {
		name := namer(e.(*asmEncoder), args[0].(data.Local))
		kwd := args[1].(data.Keyword)
		cellType, ok := cellTypes[kwd]
		if !ok {
			panic(fmt.Errorf(ErrUnknownLocalType, kwd, cellTypeNames))
		}
		e.AddLocal(name, cellType)
	}
}

func publicNamer(_ *asmEncoder, l data.Local) data.Local {
	return l
}

func privateNamer(e *asmEncoder, l data.Local) data.Local {
	p := gen.Local(l)
	e.private[l] = p
	return p
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
		if isa.IsValidOperand(int(val)) {
			return isa.Operand(val), nil
		}
	}
	return 0, fmt.Errorf(isa.ErrExpectedOperand, val)
}

func makeCellTypeNames() string {
	res := basics.MapKeys(cellTypes)
	var buf bytes.Buffer
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
