package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

type arrayWrapper struct {
	typ  reflect.Type
	len  int
	elem Wrapper
}

func makeWrappedArray(t reflect.Type) (Wrapper, error) {
	w, err := wrapType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &arrayWrapper{
		typ:  t,
		len:  t.Len(),
		elem: w,
	}, nil
}

func (a *arrayWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	vLen := v.Len()
	out := make(data.Values, vLen)
	for i := 0; i < vLen; i++ {
		elem, err := a.elem.Wrap(c, v.Index(i))
		if err != nil {
			return nil, err
		}
		out[i] = elem
	}
	return data.NewVector(out...), nil
}

func (a *arrayWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	in := sequence.ToValues(v.(data.Sequence))
	inLen := len(in)
	out := reflect.New(a.typ).Elem()
	for i := 0; i < inLen; i++ {
		v, err := a.elem.Unwrap(in[i])
		if err != nil {
			return _emptyValue, err
		}
		out.Index(i).Set(v)
	}
	return out, nil
}
