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

func (a *arrayWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	if !v.IsValid() {
		return data.Nil, nil
	}
	vLen := v.Len()
	out := make(data.Vector, vLen)
	for i := 0; i < vLen; i++ {
		elem, err := a.elem.Wrap(c, v.Index(i))
		if err != nil {
			return nil, err
		}
		out[i] = elem
	}
	return out, nil
}

func (a *arrayWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	in := sequence.ToValues(v.(data.Sequence))
	inLen := len(in)
	out := reflect.New(a.typ).Elem()
	for i := 0; i < inLen; i++ {
		v, err := a.elem.Unwrap(in[i])
		if err != nil {
			return emptyReflectValue, err
		}
		out.Index(i).Set(v)
	}
	return out, nil
}

func makeWrappedSlice(t reflect.Type) Wrapper {
	return &sliceWrapper{
		typ:  t,
		elem: wrapType(t.Elem()),
	}
}

func (s *sliceWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	if !v.IsValid() {
		return data.Nil, nil
	}
	c, err := c.Push(v)
	if err != nil {
		return nil, err
	}
	vLen := v.Len()
	out := make(data.Vector, vLen)
	for i := 0; i < vLen; i++ {
		v, err := s.elem.Wrap(c, v.Index(i))
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

func (s *sliceWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	in := sequence.ToValues(v.(data.Sequence))
	inLen := len(in)
	out := reflect.MakeSlice(s.typ, inLen, inLen)
	for i := 0; i < inLen; i++ {
		v, err := s.elem.Unwrap(in[i])
		if err != nil {
			return emptyReflectValue, err
		}
		out.Index(i).Set(v)
	}
	return out, nil
}
