package generate

import (
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
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
	args := sequence.ToVector(r)
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
		if s, ok := e.ResolveScoped(l); ok {
			switch s.Scope {
			case encoder.LocalScope, encoder.ArgScope:
				return callDynamic(e, l, args)
			case encoder.ClosureScope:
				if isSelfCalling(e, s) {
					return callSelf(e, args)
				}
			}
		}
	}
	globals := e.Globals()
	if v, err := env.ResolveValue(globals, s); err == nil {
		switch v := v.(type) {
		case compiler.Call:
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
	switch v := v.(type) {
	case compiler.Call:
		return v(e, args...)
	case data.Procedure:
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
	e.Emit(isa.Call, isa.Operand(al))
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

func callSelf(e encoder.Encoder, args data.Vector) error {
	al, err := makeArgs(e, args)()
	if err != nil {
		return err
	}
	e.Emit(isa.CallSelf, isa.Operand(al))
	return nil
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

func isSelfCalling(e encoder.Encoder, s *encoder.ScopedCell) bool {
	path := makeEncoderPath(e)
	for len(path) > 0 { // walk up to the nearest procEncoder
		pe := path[0]
		path = path[1:]
		if _, ok := pe.(*procEncoder); ok {
			break
		}
	}
	if len(path) > 0 { // we should find our bindEncoder after that
		if b, ok := path[0].(*bindEncoder); ok {
			return b.cell.Name == s.Cell.Name
		}
	}
	return false
}

func makeEncoderPath(e encoder.Encoder) []encoder.Encoder {
	res := []encoder.Encoder{e}
	last := e
	for {
		w, ok := e.(encoder.WrappedEncoder)
		if !ok {
			break
		}
		if e = w.Wrapped(); e == nil {
			break
		}
		if e != last {
			res = append(res, e)
			last = e
		}
	}
	return res
}
