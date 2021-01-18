package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

type (
	arrayWrapper struct {
		typ  reflect.Type
		len  int
		elem Wrapper
	}

	sliceWrapper struct {
		typ  reflect.Type
		elem Wrapper
	}
)

func makeWrappedArray(t reflect.Type) Wrapper {
	return &arrayWrapper{
		typ:  t,
		len:  t.Len(),
		elem: wrapType(t.Elem()),
	}
}

func (a *arrayWrapper) Wrap(v reflect.Value) data.Value {
	vLen := v.Len()
	out := make(data.Vector, vLen)
	for i := 0; i < vLen; i++ {
		out[i] = a.elem.Wrap(v.Index(i))
	}
	return out
}

func (a *arrayWrapper) Unwrap(v data.Value) reflect.Value {
	in := sequence.ToValues(v.(data.Sequence))
	inLen := len(in)
	out := reflect.New(a.typ).Elem()
	for i := 0; i < inLen; i++ {
		v := a.elem.Unwrap(in[i])
		out.Index(i).Set(v)
	}
	return out
}

func makeWrappedSlice(t reflect.Type) Wrapper {
	return &sliceWrapper{
		typ:  t,
		elem: wrapType(t.Elem()),
	}
}

func (s *sliceWrapper) Wrap(v reflect.Value) data.Value {
	vLen := v.Len()
	out := make(data.Vector, vLen)
	for i := 0; i < vLen; i++ {
		out[i] = s.elem.Wrap(v.Index(i))
	}
	return out
}

func (s *sliceWrapper) Unwrap(v data.Value) reflect.Value {
	in := sequence.ToValues(v.(data.Sequence))
	inLen := len(in)
	out := reflect.MakeSlice(s.typ, inLen, inLen)
	for i := 0; i < inLen; i++ {
		v := s.elem.Unwrap(in[i])
		out.Index(i).Set(v)
	}
	return out
}
