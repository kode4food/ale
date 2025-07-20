package ffi

import (
	"reflect"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

type inOutWrappers struct {
	in  Wrappers
	out Wrappers
}

func (io *inOutWrappers) wrap(t reflect.Type) (err error) {
	cIn := t.NumIn()
	io.in = make(Wrappers, cIn)
	for i := range cIn {
		if io.in[i], err = WrapType(t.In(i)); err != nil {
			return
		}
	}
	cOut := t.NumOut()
	io.out = make(Wrappers, cOut)
	for i := range cOut {
		if io.out[i], err = WrapType(t.Out(i)); err != nil {
			return
		}
	}
	return nil
}

func (io *inOutWrappers) wrapFunction(fn reflect.Value) data.Procedure {
	switch len(io.out) {
	case 0:
		return io.wrapVoidFunction(fn)
	case 1:
		return io.wrapValueFunction(fn)
	default:
		return io.wrapVectorFunction(fn)
	}
}

func (io *inOutWrappers) wrapValueFunction(fn reflect.Value) data.Procedure {
	return data.MakeProcedure(func(args ...ale.Value) ale.Value {
		in := io.in.mustUnwrap(args)
		out := fn.Call(in)
		res, err := io.out[0].Wrap(new(Context), out[0])
		if err != nil {
			panic(err)
		}
		return res
	}, len(io.in))
}

func (io *inOutWrappers) wrapVoidFunction(fn reflect.Value) data.Procedure {
	return data.MakeProcedure(func(args ...ale.Value) ale.Value {
		in := io.in.mustUnwrap(args)
		fn.Call(in)
		return data.Null
	}, len(io.in))
}

func (io *inOutWrappers) wrapVectorFunction(fn reflect.Value) data.Procedure {
	return data.MakeProcedure(func(args ...ale.Value) ale.Value {
		in := io.in.mustUnwrap(args)
		res := fn.Call(in)
		return io.out.mustWrap(res)
	}, len(io.in))
}
