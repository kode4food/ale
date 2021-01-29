package ffi

import (
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

const wrappedFuncType = data.Name("wrapped-func")

func makeWrappedFunc(t reflect.Type) Wrapper {
	cIn := t.NumIn()
	in := make([]Wrapper, cIn)
	for i := 0; i < cIn; i++ {
		in[i] = wrapType(t.In(i))
	}
	cOut := t.NumOut()
	out := make([]Wrapper, cOut)
	for i := 0; i < cOut; i++ {
		out[i] = wrapType(t.Out(i))
	}
	return &funcWrapper{
		typ: t,
		in:  in,
		out: out,
	}
}

func (f *funcWrapper) Wrap(_ *WrapContext, v reflect.Value) data.Value {
	if !v.IsValid() {
		return data.Nil
	}
	switch len(f.out) {
	case 0:
		return f.wrapVoidFunction(v)
	case 1:
		return f.wrapValueFunction(v)
	default:
		return f.wrapVectorFunction(v)
	}
}

func (f *funcWrapper) wrapVoidFunction(fn reflect.Value) data.Function {
	inLen := len(f.in)

	return data.Applicative(func(in ...data.Value) data.Value {
		uc := &UnwrapContext{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			wIn[i] = f.in[i].Unwrap(uc, in[i])
		}
		fn.Call(wIn)
		return data.Nil
	}, inLen)
}

func (f *funcWrapper) wrapValueFunction(fn reflect.Value) data.Function {
	inLen := len(f.in)

	return data.Applicative(func(in ...data.Value) data.Value {
		uc := &UnwrapContext{}
		wc := &WrapContext{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			wIn[i] = f.in[i].Unwrap(uc, in[i])
		}
		wOut := fn.Call(wIn)
		return f.out[0].Wrap(wc, wOut[0])
	}, inLen)
}

func (f *funcWrapper) wrapVectorFunction(fn reflect.Value) data.Function {
	inLen := len(f.in)
	outLen := len(f.out)

	return data.Applicative(func(in ...data.Value) data.Value {
		uc := &UnwrapContext{}
		wc := &WrapContext{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			wIn[i] = f.in[i].Unwrap(uc, in[i])
		}
		wOut := fn.Call(wIn)
		out := make(data.Vector, outLen)
		for i := 0; i < outLen; i++ {
			out[i] = f.out[i].Wrap(wc, wOut[i])
		}
		return out
	}, inLen)
}

func (f *funcWrapper) Unwrap(_ *UnwrapContext, v data.Value) reflect.Value {
	switch v := v.(type) {
	case data.Function:
		return f.unwrapCall(v)
	default:
		return reflect.ValueOf(v)
	}
}

func (f *funcWrapper) unwrapCall(c data.Function) reflect.Value {
	var unwrapped makeFuncType
	switch len(f.out) {
	case 0:
		unwrapped = f.unwrapVoidCall(c)
	case 1:
		unwrapped = f.unwrapValueCall(c)
	default:
		unwrapped = f.unwrapVectorCall(c)
	}
	return reflect.MakeFunc(f.typ, unwrapped)
}

func (f *funcWrapper) unwrapVoidCall(c data.Function) makeFuncType {
	inLen := len(f.in)

	return func(args []reflect.Value) []reflect.Value {
		wc := &WrapContext{}
		in := make([]data.Value, len(args))
		for i := 0; i < inLen; i++ {
			in[i] = f.in[i].Wrap(wc, args[i])
		}
		c.Call(in...)
		return []reflect.Value{}
	}
}

func (f *funcWrapper) unwrapValueCall(c data.Function) makeFuncType {
	inLen := len(f.in)

	return func(args []reflect.Value) []reflect.Value {
		wc := &WrapContext{}
		uc := &UnwrapContext{}
		in := make([]data.Value, len(args))
		for i := 0; i < inLen; i++ {
			in[i] = f.in[i].Wrap(wc, args[i])
		}
		return []reflect.Value{f.out[0].Unwrap(uc, c.Call(in...))}
	}
}

func (f *funcWrapper) unwrapVectorCall(c data.Function) makeFuncType {
	inLen := len(f.in)
	outLen := len(f.out)

	return func(args []reflect.Value) []reflect.Value {
		wc := &WrapContext{}
		uc := &UnwrapContext{}
		in := make([]data.Value, len(args))
		for i := 0; i < inLen; i++ {
			in[i] = f.in[i].Wrap(wc, args[i])
		}
		res := c.Call(in...).(data.Vector)
		out := make([]reflect.Value, outLen)
		for i := 0; i < outLen; i++ {
			out[i] = f.out[i].Unwrap(uc, res[i])
		}
		return out
	}
}
