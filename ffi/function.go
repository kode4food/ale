package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	funcWrapper struct {
		typ reflect.Type
		in  []Wrapper
		out []Wrapper
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
	in := make([]Wrapper, cIn)
	for i := 0; i < cIn; i++ {
		w, err := wrapType(t.In(i))
		if err != nil {
			return nil, err
		}
		in[i] = w
	}
	cOut := t.NumOut()
	out := make([]Wrapper, cOut)
	for i := 0; i < cOut; i++ {
		w, err := wrapType(t.Out(i))
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
	inLen := len(w.in)

	return data.Applicative(func(in ...data.Value) data.Value {
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			arg, err := w.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = arg
		}
		fn.Call(wIn)
		return data.Nil
	}, inLen)
}

func (w *funcWrapper) wrapValueFunction(fn reflect.Value) data.Function {
	inLen := len(w.in)

	return data.Applicative(func(in ...data.Value) data.Value {
		c := &Context{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			arg, err := w.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = arg
		}
		wOut := fn.Call(wIn)
		res, err := w.out[0].Wrap(c, wOut[0])
		if err != nil {
			panic(err)
		}
		return res
	}, inLen)
}

func (w *funcWrapper) wrapVectorFunction(fn reflect.Value) data.Function {
	inLen := len(w.in)
	outLen := len(w.out)

	return data.Applicative(func(in ...data.Value) data.Value {
		wc := &Context{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			arg, err := w.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = arg
		}
		wOut := fn.Call(wIn)
		out := make(data.Values, outLen)
		for i := 0; i < outLen; i++ {
			res, err := w.out[i].Wrap(wc, wOut[i])
			if err != nil {
				panic(err)
			}
			out[i] = res
		}
		return data.NewVector(out...)
	}, inLen)
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
	inLen := len(w.in)

	return func(args []reflect.Value) []reflect.Value {
		wc := &Context{}
		in := make([]data.Value, len(args))
		for i := 0; i < inLen; i++ {
			arg, err := w.in[i].Wrap(wc, args[i])
			if err != nil {
				panic(err)
			}
			in[i] = arg
		}
		c.Call(in...)
		return []reflect.Value{}
	}
}

func (w *funcWrapper) unwrapValueCall(c data.Function) makeFuncType {
	inLen := len(w.in)

	return func(args []reflect.Value) []reflect.Value {
		wc := &Context{}
		in := make([]data.Value, len(args))
		for i := 0; i < inLen; i++ {
			arg, err := w.in[i].Wrap(wc, args[i])
			if err != nil {
				panic(err)
			}
			in[i] = arg
		}
		res, err := w.out[0].Unwrap(c.Call(in...))
		if err != nil {
			panic(err)
		}
		return []reflect.Value{res}
	}
}

func (w *funcWrapper) unwrapVectorCall(c data.Function) makeFuncType {
	inLen := len(w.in)
	outLen := len(w.out)

	return func(args []reflect.Value) []reflect.Value {
		wc := &Context{}
		in := make([]data.Value, len(args))
		for i := 0; i < inLen; i++ {
			arg, err := w.in[i].Wrap(wc, args[i])
			if err != nil {
				panic(err)
			}
			in[i] = arg
		}
		res := c.Call(in...).(data.Vector).Values()
		out := make([]reflect.Value, outLen)
		for i := 0; i < outLen; i++ {
			res, err := w.out[i].Unwrap(res[i])
			if err != nil {
				panic(err)
			}
			out[i] = res
		}
		return out
	}
}
