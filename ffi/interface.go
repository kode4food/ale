package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	interfaceWrapper struct {
		methods []*methodWrapper
	}

	methodWrapper struct {
		name string
		in   []Wrapper
		out  []Wrapper
	}
)

func makeWrappedInterface(t reflect.Type) Wrapper {
	mLen := t.NumMethod()
	methods := make([]*methodWrapper, mLen)
	for i := 0; i < mLen; i++ {
		methods[i] = makeWrappedMethod(t.Method(i))
	}
	return &interfaceWrapper{
		methods: methods,
	}
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

func (i *interfaceWrapper) Wrap(v reflect.Value) data.Value {
	res := make(data.Object, len(i.methods))
	for _, m := range i.methods {
		k := data.Keyword(m.name)
		res[k] = m.wrapMethod(v)
	}
	return res
}

func (m *methodWrapper) wrapMethod(v reflect.Value) data.Call {
	switch len(m.out) {
	case 0:
		return m.wrapVoidMethod(v)
	case 1:
		return m.wrapValueMethod(v)
	default:
		return m.wrapVectorMethod(v)
	}
}

func (m *methodWrapper) wrapVoidMethod(v reflect.Value) data.Call {
	inLen := len(m.in)
	fn := v.MethodByName(m.name)

	return func(in ...data.Value) data.Value {
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			wIn[i] = m.in[i].Unwrap(in[i])
		}
		fn.Call(wIn)
		return data.Nil
	}
}

func (m *methodWrapper) wrapValueMethod(v reflect.Value) data.Call {
	inLen := len(m.in)
	fn := v.MethodByName(m.name)

	return func(in ...data.Value) data.Value {
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			wIn[i] = m.in[i].Unwrap(in[i])
		}
		wOut := fn.Call(wIn)
		return m.out[0].Wrap(wOut[0])
	}
}

func (m *methodWrapper) wrapVectorMethod(v reflect.Value) data.Call {
	inLen := len(m.in)
	outLen := len(m.out)
	fn := v.MethodByName(m.name)

	return func(in ...data.Value) data.Value {
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			wIn[i] = m.in[i].Unwrap(in[i])
		}
		wOut := fn.Call(wIn)
		out := make(data.Vector, outLen)
		for i := 0; i < outLen; i++ {
			out[i] = m.out[i].Wrap(wOut[i])
		}
		return out
	}
}

func (i *interfaceWrapper) Unwrap(data.Value) reflect.Value {
	panic("not implemented")
}
