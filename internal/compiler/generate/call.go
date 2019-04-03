package generate

import (
	"fmt"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/compiler"
	"gitlab.com/kode4food/ale/internal/compiler/encoder"
	"gitlab.com/kode4food/ale/internal/namespace"
	"gitlab.com/kode4food/ale/internal/runtime/isa"
	"gitlab.com/kode4food/ale/stdlib"
)

type argsGen func(encoder.Type, api.Values)

// Error messages
const (
	UnknownConvention = "unknown calling convention: %s"
)

// Call encodes a function call
func Call(e encoder.Type, l *api.List) {
	if l.Count() == 0 {
		Literal(e, api.EmptyList)
		return
	}
	f := l.First()
	args := stdlib.SequenceToValues(l.Rest())
	if s, ok := f.(api.Symbol); ok {
		callSymbol(e, s, args)
		return
	}
	if c, ok := f.(api.Caller); ok && !compiler.IsEvaluable(f) {
		callApplicative(e, c.Caller(), args)
		return
	}
	callDynamic(e, f, args)
}

func callSymbol(e encoder.Type, s api.Symbol, args api.Values) {
	if l, ok := s.(api.LocalSymbol); ok {
		if _, ok := e.ResolveLocal(l); ok {
			callDynamic(e, l, args)
			return
		}
	}
	globals := e.Globals()
	if v, ok := namespace.ResolveSymbol(globals, s); ok {
		switch typed := v.(type) {
		case api.Call:
			callApplicative(e, typed, args)
			return
		case *api.Function:
			if typed.IsBound() {
				callFunction(e, typed, args)
				return
			}
		}
	}
	callDynamic(e, s, args)
}

func callFunction(e encoder.Type, f *api.Function, args api.Values) {
	al := len(args)
	if err := f.CheckArity(al); err != nil {
		panic(err)
	}
	switch f.Convention {
	case api.SpecialCall:
		callSpecial(e, f.Call, args)
	case api.ApplicativeCall:
		callApplicative(e, f.Call, args)
	case api.NormalCall, api.MacroCall:
		callNormal(e, f.Call, args)
	default:
		c := f.Convention
		panic(fmt.Sprintf(UnknownConvention, c))
	}
}

func callSpecial(e encoder.Type, generate api.Call, args api.Values) {
	specialArgs := append(api.Values{e}, args...)
	generate(specialArgs...)
}

func callDynamic(e encoder.Type, v api.Value, args api.Values) {
	al := len(args)
	switch al {
	case 0:
		Value(e, v)
		e.Append(isa.MakeCall, isa.Call0)
	case 1:
		applicativeArgs(e, args)
		Value(e, v)
		e.Append(isa.MakeCall, isa.Call1)
	default:
		applicativeArgs(e, args)
		Value(e, v)
		e.Append(isa.MakeCall, isa.Call, isa.Count(al))
	}
}

func callApplicative(e encoder.Type, f api.Call, args api.Values) {
	callWith(applicativeArgs, e, f, args)
}

func callNormal(e encoder.Type, f api.Call, args api.Values) {
	callWith(normalArgs, e, f, args)
}

func callWith(gen argsGen, e encoder.Type, f api.Call, args api.Values) {
	al := len(args)
	switch al {
	case 0:
		Literal(e, f)
		e.Append(isa.Call0)
	case 1:
		gen(e, args)
		Literal(e, f)
		e.Append(isa.Call1)
	default:
		gen(e, args)
		Literal(e, f)
		e.Append(isa.Call, isa.Count(al))
	}
}

func applicativeArgs(e encoder.Type, args api.Values) {
	for i := len(args) - 1; i >= 0; i-- {
		Value(e, args[i])
	}
}

func normalArgs(e encoder.Type, args api.Values) {
	for i := len(args) - 1; i >= 0; i-- {
		Literal(e, args[i])
	}
}
