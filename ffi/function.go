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
	switch len(w.out) {
	case 0:
		return w.wrapVoidFunction(v), nil
	case 1:
		return w.wrapValueFunction(v), nil
	default:
		return w.wrapVectorFunction(v), nil
	}
}

func (w *funcWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if v, ok := v.(data.Function); ok {
		return w.unwrapCall(v), nil
	}
	return _emptyValue, errors.New(ErrValueMustBeFunction)
}

func (w *funcWrapper) unwrapCall(c data.Function) reflect.Value {
	var unwrapped makeFuncType
	switch len(w.out) {
	case 0:
		unwrapped = w.unwrapVoidCall(c)
	case 1:
		unwrapped = w.unwrapValueCall(c)
	default:
		unwrapped = w.unwrapVectorCall(c)
	}
	return reflect.MakeFunc(w.typ, unwrapped)
}

func (w *funcWrapper) unwrapVoidCall(c data.Function) makeFuncType {
	return func(args []reflect.Value) []reflect.Value {
		in := w.in.wrap(args)
		c.Call(in...)
		return []reflect.Value{}
	}
}

func (w *funcWrapper) unwrapValueCall(c data.Function) makeFuncType {
	return func(args []reflect.Value) []reflect.Value {
		in := w.in.wrap(args)
		res, err := w.out[0].Unwrap(c.Call(in...))
		if err != nil {
			panic(err)
		}
		return []reflect.Value{res}
	}
}

func (w *funcWrapper) unwrapVectorCall(c data.Function) makeFuncType {
	return func(args []reflect.Value) []reflect.Value {
		in := w.in.wrap(args)
		res := c.Call(in...).(data.Vector).Values()
		return w.out.unwrap(res)
	}
}
