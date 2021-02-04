package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

type sliceWrapper struct {
	typ  reflect.Type
	elem Wrapper
}

func makeWrappedSlice(t reflect.Type) (Wrapper, error) {
	w, err := wrapType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &sliceWrapper{
		typ:  t,
		elem: w,
	}, nil
}

func (s *sliceWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
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
			return _emptyValue, err
		}
		out.Index(i).Set(v)
	}
	return out, nil
}
