package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	interfaceWrapper []*methodWrapper

	methodWrapper struct {
		name string
		in   []Wrapper
		out  []Wrapper
	}
)

func makeWrappedInterface(t reflect.Type) Wrapper {
	mLen := t.NumMethod()
	res := make(interfaceWrapper, 0, mLen)
	for i := 0; i < mLen; i++ {
		m := t.Method(i)
		if m.PkgPath != "" { // Not exported
			continue
		}
		res = append(res, makeWrappedMethod(m))
	}
	return res
}

func makeWrappedMethod(m reflect.Method) *methodWrapper {
	t := m.Type
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
	return &methodWrapper{
		name: m.Name,
		in:   in,
		out:  out,
	}
}

func (i interfaceWrapper) Wrap(_ *WrapContext, v reflect.Value) data.Value {
	if !v.IsValid() {
		return data.Nil
	}
	res := make(data.Object, len(i))
	for _, m := range i {
		k := data.Keyword(m.name)
		res[k] = m.wrapMethod(v)
	}
	return res
}

func (m *methodWrapper) wrapMethod(v reflect.Value) data.Function {
	switch len(m.out) {
	case 0:
		return m.wrapVoidMethod(v)
	case 1:
		return m.wrapValueMethod(v)
	default:
		return m.wrapVectorMethod(v)
	}
}

func (m *methodWrapper) wrapVoidMethod(v reflect.Value) data.Function {
	inLen := len(m.in)
	fn := v.MethodByName(m.name)

	return data.Applicative(func(in ...data.Value) data.Value {
		uc := &UnwrapContext{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			wIn[i] = m.in[i].Unwrap(uc, in[i])
		}
		fn.Call(wIn)
		return data.Nil
	}, inLen)
}

func (m *methodWrapper) wrapValueMethod(v reflect.Value) data.Function {
	inLen := len(m.in)
	fn := v.MethodByName(m.name)

	return data.Applicative(func(in ...data.Value) data.Value {
		wc := &WrapContext{}
		uc := &UnwrapContext{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			wIn[i] = m.in[i].Unwrap(uc, in[i])
		}
		wOut := fn.Call(wIn)
		return m.out[0].Wrap(wc, wOut[0])
	}, inLen)
}

func (m *methodWrapper) wrapVectorMethod(v reflect.Value) data.Function {
	inLen := len(m.in)
	outLen := len(m.out)
	fn := v.MethodByName(m.name)

	return data.Applicative(func(in ...data.Value) data.Value {
		wc := &WrapContext{}
		uc := &UnwrapContext{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			wIn[i] = m.in[i].Unwrap(uc, in[i])
		}
		wOut := fn.Call(wIn)
		out := make(data.Vector, outLen)
		for i := 0; i < outLen; i++ {
			out[i] = m.out[i].Wrap(wc, wOut[i])
		}
		return out
	}, inLen)
}

func (i interfaceWrapper) Unwrap(_ *UnwrapContext, _ data.Value) reflect.Value {
	panic("not implemented")
}
