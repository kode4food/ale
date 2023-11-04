package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	funcWrapper struct {
		*inOutWrappers
		typ reflect.Type
	}

	// the type accepted by reflect.MakeFunc
	makeFuncType func(args []reflect.Value) (results []reflect.Value)
)

// Error messages
const (
	ErrValueMustBeFunction = "value must be a function"
)

func makeWrappedFunc(t reflect.Type) (Wrapper, error) {
	io, err := makeInOutWrappers(t)
	if err != nil {
		return nil, err
	}
	return &funcWrapper{
		typ:           t,
		inOutWrappers: io,
	}, nil
}

func (w *funcWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return w.wrapFunction(v), nil
}

func (w *funcWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if v, ok := v.(data.Procedure); ok {
		return w.unwrapCall(v), nil
	}
	return _emptyValue, errors.New(ErrValueMustBeFunction)
}

func (w *funcWrapper) unwrapCall(l data.Procedure) reflect.Value {
	var unwrapped makeFuncType
	switch len(w.out) {
	case 0:
		unwrapped = w.unwrapVoidCall(l)
	case 1:
		unwrapped = w.unwrapValueCall(l)
	default:
		unwrapped = w.unwrapVectorCall(l)
	}
	return reflect.MakeFunc(w.typ, unwrapped)
}

func (w *funcWrapper) unwrapVoidCall(l data.Procedure) makeFuncType {
	return func(args []reflect.Value) []reflect.Value {
		in := w.in.wrap(args)
		l.Call(in...)
		return []reflect.Value{}
	}
}

func (w *funcWrapper) unwrapValueCall(l data.Procedure) makeFuncType {
	return func(args []reflect.Value) []reflect.Value {
		in := w.in.wrap(args)
		res, err := w.out[0].Unwrap(l.Call(in...))
		if err != nil {
			panic(err)
		}
		return []reflect.Value{res}
	}
}

func (w *funcWrapper) unwrapVectorCall(l data.Procedure) makeFuncType {
	return func(args []reflect.Value) []reflect.Value {
		in := w.in.wrap(args)
		res := l.Call(in...).(data.Vector).Values()
		return w.out.unwrap(res)
	}
}
