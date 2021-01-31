package ffi

import (
	"errors"
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

// Error messages
const (
	ErrNotImplemented = "not implemented"
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

func (i interfaceWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	if !v.IsValid() {
		return data.Nil, nil
	}
	res := make(data.Object, len(i))
	for _, m := range i {
		k := data.Keyword(m.name)
		res[k] = m.wrapMethod(v)
	}
	return res, nil
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
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			u, err := m.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = u
		}
		fn.Call(wIn)
		return data.Nil
	}, inLen)
}

func (m *methodWrapper) wrapValueMethod(v reflect.Value) data.Function {
	inLen := len(m.in)
	fn := v.MethodByName(m.name)

	return data.Applicative(func(in ...data.Value) data.Value {
		c := &Context{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			arg, err := m.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = arg
		}
		wOut := fn.Call(wIn)
		res, err := m.out[0].Wrap(c, wOut[0])
		if err != nil {
			panic(err)
		}
		return res
	}, inLen)
}

func (m *methodWrapper) wrapVectorMethod(v reflect.Value) data.Function {
	inLen := len(m.in)
	outLen := len(m.out)
	fn := v.MethodByName(m.name)

	return data.Applicative(func(in ...data.Value) data.Value {
		c := &Context{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			arg, err := m.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = arg
		}
		wOut := fn.Call(wIn)
		out := make(data.Vector, outLen)
		for i := 0; i < outLen; i++ {
			res, err := m.out[i].Wrap(c, wOut[i])
			if err != nil {
				panic(err)
			}
			out[i] = res
		}
		return out
	}, inLen)
}

func (i interfaceWrapper) Unwrap(_ data.Value) (reflect.Value, error) {
	return emptyReflectValue, errors.New(ErrNotImplemented)
}
