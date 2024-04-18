package generate

import (
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/compiler"
	"github.com/kode4food/ale/pkg/compiler/encoder"
	"github.com/kode4food/ale/pkg/compiler/special"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/runtime/isa"
)

type (
	funcEmitter func()
	argsEmitter func() int
)

// Call encodes a function call
func Call(e encoder.Encoder, l *data.List) {
	f, r, ok := l.Split()
	if !ok {
		Literal(e, data.Null)
		return
	}
	args := sequence.ToValues(r)
	callValue(e, f, args)
}

func callValue(e encoder.Encoder, v data.Value, args data.Vector) {
	if s, ok := v.(data.Symbol); ok {
		callSymbol(e, s, args)
		return
	}
	callNonSymbol(e, v, args)
}

func callSymbol(e encoder.Encoder, s data.Symbol, args data.Vector) {
	if l, ok := s.(data.Local); ok {
		if _, ok := e.ResolveLocal(l); ok {
			callDynamic(e, l, args)
			return
		}
	}
	globals := e.Globals()
	if v, ok := env.ResolveValue(globals, s); ok {
		switch v := v.(type) {
		case special.Call:
			v(e, args...)
			return
		case data.Procedure:
			callStatic(e, v, args)
			return
		}
	}
	callDynamic(e, s, args)
}

func callNonSymbol(e encoder.Encoder, v data.Value, args data.Vector) {
	if compiler.IsEvaluable(v) {
		callDynamic(e, v, args)
		return
	}
	switch v := v.(type) {
	case data.Procedure:
		callStatic(e, v, args)
	default:
		callDynamic(e, v, args)
	}
}

func assertArity(p data.Procedure, args data.Vector) {
	al := len(args)
	if err := p.CheckArity(al); err != nil {
		panic(err)
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
		e.Emit(isa.Call, isa.Operand(al))
	}
}

func callStatic(e encoder.Encoder, p data.Procedure, args data.Vector) {
	assertArity(p, args)
	emitFunc := staticLiteral(e, p)
	emitArgs := makeArgs(e, args)
	callWith(e, emitFunc, emitArgs)
}

func staticLiteral(e encoder.Encoder, fn data.Value) funcEmitter {
	return func() {
		Literal(e, fn)
	}
}

func callDynamic(e encoder.Encoder, v data.Value, args data.Vector) {
	emitFunc := dynamicEval(e, v)
	emitArgs := makeArgs(e, args)
	callWith(e, emitFunc, emitArgs)
}

func dynamicEval(e encoder.Encoder, v data.Value) funcEmitter {
	return func() {
		Value(e, v)
	}
}

func makeArgs(e encoder.Encoder, args data.Vector) argsEmitter {
	return func() int {
		al := len(args)
		for i := al - 1; i >= 0; i-- {
			Value(e, args[i])
		}
		return al
	}
}
