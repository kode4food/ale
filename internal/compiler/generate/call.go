package generate

import (
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/special"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

type (
	funcEmitter func() error
	argsEmitter func() (int, error)
)

// Call encodes a function call
func Call(e encoder.Encoder, l *data.List) error {
	f, r, ok := l.Split()
	if !ok {
		return Literal(e, data.Null)
	}
	args := sequence.ToValues(r)
	return callValue(e, f, args)
}

func callValue(e encoder.Encoder, v data.Value, args data.Vector) error {
	if s, ok := v.(data.Symbol); ok {
		return callSymbol(e, s, args)
	}
	return callNonSymbol(e, v, args)
}

func callSymbol(e encoder.Encoder, s data.Symbol, args data.Vector) error {
	if l, ok := s.(data.Local); ok {
		if _, ok := e.ResolveLocal(l); ok {
			return callDynamic(e, l, args)
		}
	}
	globals := e.Globals()
	if v, err := env.ResolveValue(globals, s); err == nil {
		switch v := v.(type) {
		case special.Call:
			return v(e, args...)
		case data.Procedure:
			return callStatic(e, v, args)
		}
	}
	return callDynamic(e, s, args)
}

func callNonSymbol(e encoder.Encoder, v data.Value, args data.Vector) error {
	if compiler.IsEvaluable(v) {
		return callDynamic(e, v, args)
	}
	if v, ok := v.(data.Procedure); ok {
		return callStatic(e, v, args)
	}
	return callDynamic(e, v, args)
}

func callWith(e encoder.Encoder, fn funcEmitter, args argsEmitter) error {
	al, err := args()
	if err != nil {
		return err
	}
	if err := fn(); err != nil {
		return err
	}
	switch al {
	case 0:
		e.Emit(isa.Call0)
	case 1:
		e.Emit(isa.Call1)
	default:
		e.Emit(isa.Call, isa.Operand(al))
	}
	return nil
}

func callStatic(e encoder.Encoder, p data.Procedure, args data.Vector) error {
	if err := p.CheckArity(len(args)); err != nil {
		return err
	}
	emitFunc := staticLiteral(e, p)
	emitArgs := makeArgs(e, args)
	return callWith(e, emitFunc, emitArgs)
}

func staticLiteral(e encoder.Encoder, fn data.Value) funcEmitter {
	return func() error {
		return Literal(e, fn)
	}
}

func callDynamic(e encoder.Encoder, v data.Value, args data.Vector) error {
	emitFunc := dynamicEval(e, v)
	emitArgs := makeArgs(e, args)
	return callWith(e, emitFunc, emitArgs)
}

func dynamicEval(e encoder.Encoder, v data.Value) funcEmitter {
	return func() error {
		return Value(e, v)
	}
}

func makeArgs(e encoder.Encoder, args data.Vector) argsEmitter {
	return func() (int, error) {
		al := len(args)
		for i := al - 1; i >= 0; i-- {
			if err := Value(e, args[i]); err != nil {
				return 0, err
			}
		}
		return al, nil
	}
}
