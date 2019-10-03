package generate

import (
	"fmt"

	"github.com/kode4food/ale/compiler"
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/namespace"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/ale/runtime/vm"
	"github.com/kode4food/ale/stdlib"
)

type (
	funcEmitter func()
	argsEmitter func() int
	valEmitter  func(encoder.Type, data.Value)
)

// Error messages
const (
	UnknownConvention = "unknown calling convention: %s"
)

// Call encodes a function call
func Call(e encoder.Type, l data.List) {
	if l.Count() == 0 {
		Literal(e, data.EmptyList)
		return
	}
	f := l.First()
	args := stdlib.SequenceToValues(l.Rest())
	callValue(e, f, args)
}

func callValue(e encoder.Type, v data.Value, args data.Values) {
	if s, ok := v.(data.Symbol); ok {
		callSymbol(e, s, args)
		return
	}
	callNonSymbol(e, v, args)
}

func callSymbol(e encoder.Type, s data.Symbol, args data.Values) {
	if l, ok := s.(data.LocalSymbol); ok {
		if _, ok := e.ResolveLocal(l.Name()); ok {
			callDynamic(e, l, args)
			return
		}
	}
	globals := e.Globals()
	if v, ok := namespace.ResolveValue(globals, s); ok {
		switch typed := v.(type) {
		case encoder.Call:
			typed(e, args...)
			return
		case data.Call:
			callApplicative(e, typed, args)
			return
		case data.Function:
			callFunction(e, typed, args)
			return
		}
	}
	callDynamic(e, s, args)
}

func callNonSymbol(e encoder.Type, v data.Value, args data.Values) {
	if compiler.IsEvaluable(v) {
		callDynamic(e, v, args)
		return
	}
	switch typed := v.(type) {
	case data.Function:
		callFunction(e, typed, args)
	case data.Caller:
		callCaller(e, typed, args)
	default:
		callDynamic(e, typed, args)
	}
}

func callCaller(e encoder.Type, c data.Caller, args data.Values) {
	emitFunc := callerLiteral(e, c)
	emitArgs := applicativeArgs(e, args)
	callWith(e, emitFunc, emitArgs)
}

func callFunction(e encoder.Type, f data.Function, args data.Values) {
	assertArity(f, args)
	emitFunc := funcRef(e, f)
	emitArgs := funcArgs(e, f, args)
	callWith(e, emitFunc, emitArgs)
}

func assertArity(f data.Function, args data.Values) {
	al := len(args)
	if err := f.CheckArity(al); err != nil {
		panic(err)
	}
}

func funcRef(e encoder.Type, f data.Function) funcEmitter {
	if cl, ok := f.(*vm.Closure); ok {
		return dynamicLiteral(e, cl)
	}
	return callerLiteral(e, f)
}

func funcArgs(e encoder.Type, f data.Function, args data.Values) argsEmitter {
	c := f.Convention()
	switch c {
	case data.ApplicativeCall:
		return applicativeArgs(e, args)
	case data.NormalCall:
		return normalArgs(e, args)
	default:
		panic(fmt.Sprintf(UnknownConvention, c))
	}
}

func callWith(e encoder.Type, emitFunc funcEmitter, emitArgs argsEmitter) {
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

func callApplicative(e encoder.Type, f data.Call, args data.Values) {
	emitFunc := staticLiteral(e, f)
	emitArgs := applicativeArgs(e, args)
	callWith(e, emitFunc, emitArgs)
}

func callDynamic(e encoder.Type, v data.Value, args data.Values) {
	emitFunc := dynamicEval(e, v)
	emitArgs := applicativeArgs(e, args)
	callWith(e, emitFunc, emitArgs)
}

func staticLiteral(e encoder.Type, fn data.Value) funcEmitter {
	return func() {
		Literal(e, fn)
	}
}

func callerLiteral(e encoder.Type, fn data.Caller) funcEmitter {
	return func() {
		Literal(e, fn.Call())
	}
}

func dynamicLiteral(e encoder.Type, fn data.Value) funcEmitter {
	return func() {
		Literal(e, fn)
		e.Emit(isa.MakeCall)
	}
}

func dynamicEval(e encoder.Type, v data.Value) funcEmitter {
	return func() {
		Value(e, v)
		e.Emit(isa.MakeCall)
	}
}

func applicativeArgs(e encoder.Type, args data.Values) argsEmitter {
	return makeArgs(e, args, Value)
}

func normalArgs(e encoder.Type, args data.Values) argsEmitter {
	return makeArgs(e, args, Literal)
}

func makeArgs(e encoder.Type, args data.Values, emit valEmitter) argsEmitter {
	return func() int {
		al := len(args)
		for i := al - 1; i >= 0; i-- {
			emit(e, args[i])
		}
		return al
	}
}
