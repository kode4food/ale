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

func (a *arrayWrapper) Wrap(c *WrapContext, v reflect.Value) data.Value {
	if r, ok := c.Get(v); ok {
		return r
	}
	vLen := v.Len()
	out := make(data.Vector, vLen)
	c.Put(v, out)
	for i := 0; i < vLen; i++ {
		out[i] = a.elem.Wrap(c, v.Index(i))
	}
	return out
}

func (a *arrayWrapper) Unwrap(c *UnwrapContext, v data.Value) reflect.Value {
	if r, ok := c.Get(v); ok {
		return r
	}
	in := sequence.ToValues(v.(data.Sequence))
	inLen := len(in)
	out := reflect.New(a.typ).Elem()
	c.Put(v, out)
	for i := 0; i < inLen; i++ {
		v := a.elem.Unwrap(c, in[i])
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

func (s *sliceWrapper) Wrap(c *WrapContext, v reflect.Value) data.Value {
	if r, ok := c.Get(v); ok {
		return r
	}
	vLen := v.Len()
	out := make(data.Vector, vLen)
	c.Put(v, out)
	for i := 0; i < vLen; i++ {
		out[i] = s.elem.Wrap(c, v.Index(i))
	}
	return out
}

func (s *sliceWrapper) Unwrap(c *UnwrapContext, v data.Value) reflect.Value {
	if r, ok := c.Get(v); ok {
		return r
	}
	in := sequence.ToValues(v.(data.Sequence))
	inLen := len(in)
	out := reflect.MakeSlice(s.typ, inLen, inLen)
	c.Put(v, out)
	for i := 0; i < inLen; i++ {
		v := s.elem.Unwrap(c, in[i])
		out.Index(i).Set(v)
	}
	return out
}
