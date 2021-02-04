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

func (f *funcWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	switch len(f.out) {
	case 0:
		return f.wrapVoidFunction(v), nil
	case 1:
		return f.wrapValueFunction(v), nil
	default:
		return f.wrapVectorFunction(v), nil
	}
}

func (f *funcWrapper) wrapVoidFunction(fn reflect.Value) data.Function {
	inLen := len(f.in)

	return data.Applicative(func(in ...data.Value) data.Value {
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			arg, err := f.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = arg
		}
		fn.Call(wIn)
		return data.Nil
	}, inLen)
}

func (f *funcWrapper) wrapValueFunction(fn reflect.Value) data.Function {
	inLen := len(f.in)

	return data.Applicative(func(in ...data.Value) data.Value {
		c := &Context{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			arg, err := f.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = arg
		}
		wOut := fn.Call(wIn)
		res, err := f.out[0].Wrap(c, wOut[0])
		if err != nil {
			panic(err)
		}
		return res
	}, inLen)
}

func (f *funcWrapper) wrapVectorFunction(fn reflect.Value) data.Function {
	inLen := len(f.in)
	outLen := len(f.out)

	return data.Applicative(func(in ...data.Value) data.Value {
		wc := &Context{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			arg, err := f.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = arg
		}
		wOut := fn.Call(wIn)
		out := make(data.Vector, outLen)
		for i := 0; i < outLen; i++ {
			res, err := f.out[i].Wrap(wc, wOut[i])
			if err != nil {
				panic(err)
			}
			out[i] = res
		}
		return out
	}, inLen)
}

func (f *funcWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	switch v := v.(type) {
	case data.Function:
		return f.unwrapCall(v), nil
	default:
		return reflect.ValueOf(v), nil
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
		wc := &Context{}
		in := make([]data.Value, len(args))
		for i := 0; i < inLen; i++ {
			arg, err := f.in[i].Wrap(wc, args[i])
			if err != nil {
				panic(err)
			}
			in[i] = arg
		}
		c.Call(in...)
		return []reflect.Value{}
	}
}

func (f *funcWrapper) unwrapValueCall(c data.Function) makeFuncType {
	inLen := len(f.in)

	return func(args []reflect.Value) []reflect.Value {
		wc := &Context{}
		in := make([]data.Value, len(args))
		for i := 0; i < inLen; i++ {
			arg, err := f.in[i].Wrap(wc, args[i])
			if err != nil {
				panic(err)
			}
			in[i] = arg
		}
		res, err := f.out[0].Unwrap(c.Call(in...))
		if err != nil {
			panic(err)
		}
		return []reflect.Value{res}
	}
}

func (f *funcWrapper) unwrapVectorCall(c data.Function) makeFuncType {
	inLen := len(f.in)
	outLen := len(f.out)

	return func(args []reflect.Value) []reflect.Value {
		wc := &Context{}
		in := make([]data.Value, len(args))
		for i := 0; i < inLen; i++ {
			arg, err := f.in[i].Wrap(wc, args[i])
			if err != nil {
				panic(err)
			}
			in[i] = arg
		}
		res := c.Call(in...).(data.Vector)
		out := make([]reflect.Value, outLen)
		for i := 0; i < outLen; i++ {
			res, err := f.out[i].Unwrap(res[i])
			if err != nil {
				panic(err)
			}
			out[i] = res
		}
		return out
	}
}
