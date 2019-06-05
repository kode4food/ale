package generate

import (
	"fmt"

	"gitlab.com/kode4food/ale/compiler"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
	"gitlab.com/kode4food/ale/runtime/vm"
	"gitlab.com/kode4food/ale/stdlib"
)

type (
	emitFunc func()
	emitArgs func(encoder.Type, data.Values)
)

// Error messages
const (
	UnknownConvention = "unknown calling convention: %s"
)

// Call encodes a function call
func Call(e encoder.Type, l *data.List) {
	if l.Count() == 0 {
		Literal(e, data.EmptyList)
		return
	}
	f := l.First()
	args := stdlib.SequenceToValues(l.Rest())
	if s, ok := f.(data.Symbol); ok {
		callSymbol(e, s, args)
		return
	}
	if c, ok := f.(data.Caller); ok && !compiler.IsEvaluable(f) {
		callCaller(e, c, args)
		return
	}
	callDynamic(e, f, args)
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

func callFunction(e encoder.Type, f data.Function, args data.Values) {
	assertArity(f, args)
	fEmit := functionGenerator(e, f)
	aEmit := argumentEmitter(f)
	callWith(e, fEmit, aEmit, args)
}

func assertArity(f data.Function, args data.Values) {
	al := len(args)
	if err := f.CheckArity(al); err != nil {
		panic(err)
	}
}

func functionGenerator(e encoder.Type, f data.Function) emitFunc {
	if cl, ok := f.(*vm.Closure); ok {
		return dynamicLiteral(e, cl)
	}
	return callerLiteral(e, f)
}

func argumentEmitter(f data.Function) emitArgs {
	c := f.Convention()
	switch c {
	case data.ApplicativeCall:
		return applicativeArgs
	case data.NormalCall:
		return normalArgs
	default:
		panic(fmt.Sprintf(UnknownConvention, c))
	}
}

func callWith(e encoder.Type, ef emitFunc, ea emitArgs, args data.Values) {
	al := len(args)
	switch al {
	case 0:
		ef()
		e.Emit(isa.Call0)
	case 1:
		ea(e, args)
		ef()
		e.Emit(isa.Call1)
	default:
		ea(e, args)
		ef()
		e.Emit(isa.Call, isa.Count(al))
	}
}

func callApplicative(e encoder.Type, f data.Call, args data.Values) {
	callWith(e, staticLiteral(e, f), applicativeArgs, args)
}

func callCaller(e encoder.Type, c data.Caller, args data.Values) {
	callWith(e, callerLiteral(e, c), applicativeArgs, args)
}

func callDynamic(e encoder.Type, v data.Value, args data.Values) {
	callWith(e, dynamicEval(e, v), applicativeArgs, args)
}

func staticLiteral(e encoder.Type, fn data.Value) emitFunc {
	return func() {
		Literal(e, fn)
	}
}

func callerLiteral(e encoder.Type, fn data.Caller) emitFunc {
	return func() {
		Literal(e, fn.Caller())
	}
}

func dynamicLiteral(e encoder.Type, fn data.Value) emitFunc {
	return func() {
		Literal(e, fn)
		e.Emit(isa.MakeCall)
	}
}

func dynamicEval(e encoder.Type, v data.Value) emitFunc {
	return func() {
		Value(e, v)
		e.Emit(isa.MakeCall)
	}
}

func applicativeArgs(e encoder.Type, args data.Values) {
	for i := len(args) - 1; i >= 0; i-- {
		Value(e, args[i])
	}
}

func normalArgs(e encoder.Type, args data.Values) {
	for i := len(args) - 1; i >= 0; i-- {
		Literal(e, args[i])
	}
}
