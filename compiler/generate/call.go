package generate

import (
	"fmt"

	"github.com/kode4food/ale/compiler"
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/runtime/isa"
)

type (
	funcEmitter func()
	argsEmitter func() int
	valEmitter  func(encoder.Encoder, data.Value)
)

// Error messages
const (
	ErrUnknownConvention = "unknown calling convention: %s"
)

// Call encodes a function call
func Call(e encoder.Encoder, l data.List) {
	if l.Count() == 0 {
		Literal(e, data.EmptyList)
		return
	}
	f := l.First()
	args := sequence.ToValues(l.Rest())
	callValue(e, f, args)
}

func callValue(e encoder.Encoder, v data.Value, args data.Values) {
	if s, ok := v.(data.Symbol); ok {
		callSymbol(e, s, args)
		return
	}
	callNonSymbol(e, v, args)
}

func callSymbol(e encoder.Encoder, s data.Symbol, args data.Values) {
	if l, ok := s.(data.LocalSymbol); ok {
		if _, ok := e.ResolveLocal(l.Name()); ok {
			callDynamic(e, l, args)
			return
		}
	}
	globals := e.Globals()
	if v, ok := env.ResolveValue(globals, s); ok {
		switch v := v.(type) {
		case encoder.Call:
			v(e, args...)
			return
		case data.Function:
			callFunction(e, v, args)
			return
		}
	}
	callDynamic(e, s, args)
}

func callNonSymbol(e encoder.Encoder, v data.Value, args data.Values) {
	if compiler.IsEvaluable(v) {
		callDynamic(e, v, args)
		return
	}
	switch v := v.(type) {
	case data.Function:
		callFunction(e, v, args)
	default:
		callDynamic(e, v, args)
	}
}

func callFunction(e encoder.Encoder, f data.Function, args data.Values) {
	assertArity(f, args)
	emitFunc := staticLiteral(e, f)
	emitArgs := funcArgs(e, f, args)
	callWith(e, emitFunc, emitArgs)
}

func assertArity(f data.Function, args data.Values) {
	al := len(args)
	if err := f.CheckArity(al); err != nil {
		panic(err)
	}
}

func funcArgs(e encoder.Encoder, f data.Function, args data.Values) argsEmitter {
	c := f.Convention()
	switch c {
	case data.ApplicativeCall:
		return applicativeArgs(e, args)
	case data.NormalCall:
		return normalArgs(e, args)
	default:
		panic(fmt.Errorf(ErrUnknownConvention, c))
	}
}

func callWith(e encoder.Encoder, emitFunc funcEmitter, emitArgs argsEmitter) {
	al := emitArgs()
	emitFunc()
	switch al {
	case 0:
		e.Emit(isa.Call0)
	case 1:
		e.Emit(isa.Call1)
	default:
		e.Emit(isa.Call, isa.Count(al))
	}
}

func callApplicative(e encoder.Encoder, f data.Function, args data.Values) {
	emitFunc := staticLiteral(e, f)
	emitArgs := applicativeArgs(e, args)
	callWith(e, emitFunc, emitArgs)
}

func callDynamic(e encoder.Encoder, v data.Value, args data.Values) {
	emitFunc := dynamicEval(e, v)
	emitArgs := applicativeArgs(e, args)
	callWith(e, emitFunc, emitArgs)
}

func staticLiteral(e encoder.Encoder, fn data.Value) funcEmitter {
	return func() {
		Literal(e, fn)
	}
}

func dynamicEval(e encoder.Encoder, v data.Value) funcEmitter {
	return func() {
		Value(e, v)
	}
}

func applicativeArgs(e encoder.Encoder, args data.Values) argsEmitter {
	return makeArgs(e, args, Value)
}

func normalArgs(e encoder.Encoder, args data.Values) argsEmitter {
	return makeArgs(e, args, Literal)
}

func makeArgs(e encoder.Encoder, args data.Values, emit valEmitter) argsEmitter {
	return func() int {
		al := len(args)
		for i := al - 1; i >= 0; i-- {
			emit(e, args[i])
		}
		return al
	}
}
