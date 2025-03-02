package ffi

import (
	"reflect"

	"github.com/kode4food/ale/pkg/data"
)

type inOutWrappers struct {
	in  Wrappers
	out Wrappers
}

func makeInOutWrappers(t reflect.Type) (*inOutWrappers, error) {
	cIn := t.NumIn()
	in := make(Wrappers, cIn)
	for i := range cIn {
		w, err := WrapType(t.In(i))
		if err != nil {
			return nil, err
		}
		in[i] = w
	}
	cOut := t.NumOut()
	out := make(Wrappers, cOut)
	for i := range cOut {
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

func (w *inOutWrappers) wrapFunction(fn reflect.Value) data.Procedure {
	switch len(w.out) {
	case 0:
		return w.wrapVoidFunction(fn)
	case 1:
		return w.wrapValueFunction(fn)
	default:
		return w.wrapVectorFunction(fn)
	}
}

func (w *inOutWrappers) wrapValueFunction(fn reflect.Value) data.Procedure {
	return data.MakeProcedure(func(args ...data.Value) data.Value {
		in := w.in.mustUnwrap(args)
		out := fn.Call(in)
		res, err := w.out[0].Wrap(new(Context), out[0])
		if err != nil {
			panic(err)
		}
		return res
	}, len(w.in))
}

func (w *inOutWrappers) wrapVoidFunction(fn reflect.Value) data.Procedure {
	return data.MakeProcedure(func(args ...data.Value) data.Value {
		in := w.in.mustUnwrap(args)
		fn.Call(in)
		return data.Null
	}, len(w.in))
}

func (w *inOutWrappers) wrapVectorFunction(fn reflect.Value) data.Procedure {
	return data.MakeProcedure(func(args ...data.Value) data.Value {
		in := w.in.mustUnwrap(args)
		res := fn.Call(in)
		return w.out.mustWrap(res)
	}, len(w.in))
}
