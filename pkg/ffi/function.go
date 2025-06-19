package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/pkg/data"
)

type (
	funcWrapper struct {
		inOutWrappers
		typ reflect.Type
	}

	// the type accepted by reflect.MakeFunc
	makeFuncType func([]reflect.Value) []reflect.Value
)

// ErrValueMustBeProcedure is raised when a function Unwrap call can't treat
// its source as a data.Procedure
const ErrValueMustBeProcedure = "value must be a procedure"

func makeWrappedFunc(t reflect.Type) (Wrapper, error) {
	res := &funcWrapper{typ: t}
	return res, res.wrap(t)
}

func (w *funcWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return w.wrapFunction(v), nil
}

func (w *funcWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if v, ok := v.(data.Procedure); ok {
		return w.unwrapCall(v), nil
	}
	return _zero, errors.New(ErrValueMustBeProcedure)
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
		in := w.in.mustWrap(args)
		l.Call(in...)
		return []reflect.Value{}
	}
}

func (w *funcWrapper) unwrapValueCall(l data.Procedure) makeFuncType {
	return func(args []reflect.Value) []reflect.Value {
		in := w.in.mustWrap(args)
		res, err := w.out[0].Unwrap(l.Call(in...))
		if err != nil {
			panic(err)
		}
		return []reflect.Value{res}
	}
}

func (w *funcWrapper) unwrapVectorCall(l data.Procedure) makeFuncType {
	return func(args []reflect.Value) []reflect.Value {
		in := w.in.mustWrap(args)
		res := l.Call(in...).(data.Vector)
		return w.out.mustUnwrap(res)
	}
}
