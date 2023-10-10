package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	funcWrapper struct {
		typ reflect.Type
		in  Wrappers
		out Wrappers
	}

	// the type accepted by reflect.MakeFunc
	makeFuncType func(args []reflect.Value) (results []reflect.Value)
)

// Error messages
const (
	ErrValueMustBeFunction = "value must be a function"
)

func makeWrappedFunc(t reflect.Type) (Wrapper, error) {
	cIn := t.NumIn()
	in := make(Wrappers, cIn)
	for i := 0; i < cIn; i++ {
		w, err := WrapType(t.In(i))
		if err != nil {
			return nil, err
		}
		in[i] = w
	}
	cOut := t.NumOut()
	out := make(Wrappers, cOut)
	for i := 0; i < cOut; i++ {
		w, err := WrapType(t.Out(i))
		if err != nil {
			return nil, err
		}
		out[i] = w
	}
	return &funcWrapper{
		typ: t,
		in:  in,
		out: out,
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

func (w *funcWrapper) wrapVoidFunction(fn reflect.Value) data.Function {
	return data.Applicative(func(args ...data.Value) data.Value {
		in := w.in.unwrap(args)
		fn.Call(in)
		return data.Null
	}, len(w.in))
}

func (w *funcWrapper) wrapValueFunction(fn reflect.Value) data.Function {
	return data.Applicative(func(args ...data.Value) data.Value {
		in := w.in.unwrap(args)
		out := fn.Call(in)
		res, err := w.out[0].Wrap(new(Context), out[0])
		if err != nil {
			panic(err)
		}
		return res
	}, len(w.in))
}

func (w *funcWrapper) wrapVectorFunction(fn reflect.Value) data.Function {
	return data.Applicative(func(args ...data.Value) data.Value {
		in := w.in.unwrap(args)
		res := fn.Call(in)
		out := w.out.wrap(res)
		return data.NewVector(out...)
	}, len(w.in))
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
