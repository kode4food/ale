package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

type sliceWrapper struct {
	typ  reflect.Type
	elem Wrapper
}

func makeWrappedSlice(t reflect.Type) (Wrapper, error) {
	w, err := WrapType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &sliceWrapper{
		typ:  t,
		elem: w,
	}, nil
}

func (w *sliceWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	c, err := c.Push(v)
	if err != nil {
		return data.Null, err
	}
	vLen := v.Len()
	out := make(data.Values, vLen)
	for i := 0; i < vLen; i++ {
		v, err := w.elem.Wrap(c, v.Index(i))
		if err != nil {
			return data.Null, err
		}
		out[i] = v
	}
	return data.Vector(out), nil
}

func (w *sliceWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	s, ok := v.(data.Sequence)
	if !ok {
		return _emptyValue, errors.New(ErrValueMustBeSequence)
	}
	in := sequence.ToValues(s)
	inLen := len(in)
	out := reflect.MakeSlice(w.typ, inLen, inLen)
	for i, e := range in {
		v, err := w.elem.Unwrap(e)
		if err != nil {
			return _emptyValue, err
		}
		out.Index(i).Set(v)
	}
	return out, nil
}
