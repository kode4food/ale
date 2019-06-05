package generate

import (
	"fmt"

	"gitlab.com/kode4food/ale/runtime/vm"

	"gitlab.com/kode4food/ale/compiler"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
	"gitlab.com/kode4food/ale/stdlib"
)

type argsGen func(encoder.Type, data.Values)

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
		callApplicative(e, c.Caller(), args)
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
	al := len(args)
	if err := f.CheckArity(al); err != nil {
		panic(err)
	}
	if cl, ok := f.(*vm.Closure); ok {
		callDynamic(e, cl, args)
		return
	}
	c := f.Convention()
	switch c {
	case data.ApplicativeCall:
		callApplicative(e, f.Caller(), args)
	case data.NormalCall:
		callNormal(e, f.Caller(), args)
	default:
		panic(fmt.Sprintf(UnknownConvention, c))
	}
}

func callDynamic(e encoder.Type, v data.Value, args data.Values) {
	al := len(args)
	switch al {
	case 0:
		Value(e, v)
		e.Emit(isa.MakeCall)
		e.Emit(isa.Call0)
	case 1:
		applicativeArgs(e, args)
		Value(e, v)
		e.Emit(isa.MakeCall)
		e.Emit(isa.Call1)
	default:
		applicativeArgs(e, args)
		Value(e, v)
		e.Emit(isa.MakeCall)
		e.Emit(isa.Call, isa.Count(al))
	}
}

func callApplicative(e encoder.Type, f data.Call, args data.Values) {
	callWith(applicativeArgs, e, f, args)
}

func callNormal(e encoder.Type, f data.Call, args data.Values) {
	callWith(normalArgs, e, f, args)
}

func callWith(gen argsGen, e encoder.Type, f data.Call, args data.Values) {
	al := len(args)
	switch al {
	case 0:
		Literal(e, f)
		e.Emit(isa.Call0)
	case 1:
		gen(e, args)
		Literal(e, f)
		e.Emit(isa.Call1)
	default:
		gen(e, args)
		Literal(e, f)
		e.Emit(isa.Call, isa.Count(al))
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
