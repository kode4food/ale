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

	wrappedFunc struct {
		*funcWrapper
		elem reflect.Value
		call data.Call
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

func (f *funcWrapper) Wrap(v reflect.Value) data.Value {
	return &wrappedFunc{
		funcWrapper: f,
		elem:        v,
		call:        f.wrapCall(v),
	}
}

func (f *funcWrapper) wrapCall(v reflect.Value) data.Call {
	switch len(f.out) {
	case 0:
		return f.wrapVoidCall(v)
	case 1:
		return f.wrapValueCall(v)
	default:
		return f.wrapVectorCall(v)
	}
}

func (f *funcWrapper) wrapVoidCall(v reflect.Value) data.Call {
	inLen := len(f.in)

	return func(in ...data.Value) data.Value {
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			wIn[i] = f.in[i].Unwrap(in[i])
		}
		v.Call(wIn)
		return data.Nil
	}
}

func (f *funcWrapper) wrapValueCall(v reflect.Value) data.Call {
	inLen := len(f.in)

	return func(in ...data.Value) data.Value {
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			wIn[i] = f.in[i].Unwrap(in[i])
		}
		wOut := v.Call(wIn)
		return f.out[0].Wrap(wOut[0])
	}
}

func (f *funcWrapper) wrapVectorCall(v reflect.Value) data.Call {
	inLen := len(f.in)
	outLen := len(f.out)

	return func(in ...data.Value) data.Value {
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			wIn[i] = f.in[i].Unwrap(in[i])
		}
		wOut := v.Call(wIn)
		out := make(data.Vector, outLen)
		for i := 0; i < outLen; i++ {
			out[i] = f.out[i].Wrap(wOut[i])
		}
		return out
	}
}

func (f *funcWrapper) Unwrap(v data.Value) reflect.Value {
	switch v := v.(type) {
	case *wrappedFunc:
		return v.elem
	case data.Caller:
		return f.unwrapCall(v.Call())
	case data.Call:
		return f.unwrapCall(v)
	default:
		return reflect.ValueOf(v)
	}
}

func (f *funcWrapper) unwrapCall(c data.Call) reflect.Value {
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

func (f *funcWrapper) unwrapVoidCall(c data.Call) makeFuncType {
	inLen := len(f.in)

	return func(args []reflect.Value) []reflect.Value {
		in := make([]data.Value, len(args))
		for i := 0; i < inLen; i++ {
			in[i] = f.in[i].Wrap(args[i])
		}
		return []reflect.Value{f.out[0].Unwrap(c(in...))}
	}
}

func (f *funcWrapper) unwrapValueCall(c data.Call) makeFuncType {
	inLen := len(f.in)

	return func(args []reflect.Value) []reflect.Value {
		in := make([]data.Value, len(args))
		for i := 0; i < inLen; i++ {
			in[i] = f.in[i].Wrap(args[i])
		}
		return []reflect.Value{f.out[0].Unwrap(c(in...))}
	}
}

func (f *funcWrapper) unwrapVectorCall(c data.Call) makeFuncType {
	inLen := len(f.in)
	outLen := len(f.out)

	return func(args []reflect.Value) []reflect.Value {
		in := make([]data.Value, len(args))
		for i := 0; i < inLen; i++ {
			in[i] = f.in[i].Wrap(args[i])
		}
		res := c(in...).(data.Vector)
		out := make([]reflect.Value, outLen)
		for i := 0; i < outLen; i++ {
			out[i] = f.out[i].Unwrap(res[i])
		}
		return out
	}
}

func (w *wrappedFunc) Call() data.Call {
	return w.call
}

func (*wrappedFunc) Type() data.Name {
	return wrappedFuncType
}

func (w *wrappedFunc) String() string {
	return data.DumpString(w)
}
