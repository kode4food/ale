package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
)

type inOutWrappers struct {
	in  Wrappers
	out Wrappers
}

func makeInOutWrappers(t reflect.Type) (*inOutWrappers, error) {
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
	return &inOutWrappers{
		in:  in,
		out: out,
	}, nil
}

func (w *inOutWrappers) wrapFunction(fn reflect.Value) data.Lambda {
	switch len(w.out) {
	case 0:
		return w.wrapVoidFunction(fn)
	case 1:
		return w.wrapValueFunction(fn)
	default:
		return w.wrapVectorFunction(fn)
	}
}

func (w *inOutWrappers) wrapValueFunction(fn reflect.Value) data.Lambda {
	return data.MakeLambda(func(args ...data.Value) data.Value {
		in := w.in.unwrap(args)
		out := fn.Call(in)
		res, err := w.out[0].Wrap(new(Context), out[0])
		if err != nil {
			panic(err)
		}
		return res
	}, len(w.in))
}

func (w *inOutWrappers) wrapVoidFunction(fn reflect.Value) data.Lambda {
	return data.MakeLambda(func(args ...data.Value) data.Value {
		in := w.in.unwrap(args)
		fn.Call(in)
		return data.Null
	}, len(w.in))
}

func (w *inOutWrappers) wrapVectorFunction(fn reflect.Value) data.Lambda {
	return data.MakeLambda(func(args ...data.Value) data.Value {
		in := w.in.unwrap(args)
		res := fn.Call(in)
		out := w.out.wrap(res)
		return data.NewVector(out...)
	}, len(w.in))
}
