package ffi

import (
	"reflect"

	"github.com/kode4food/ale/pkg/data"
)

type pointerWrapper struct {
	elem Wrapper
}

func makeWrappedPointer(t reflect.Type) (Wrapper, error) {
	w, err := WrapType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &pointerWrapper{
		elem: w,
	}, nil
}

func (w *pointerWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	c, err := c.Push(v)
	if err != nil {
		return data.Null, err
	}
	e := v.Elem()
	if e.IsValid() {
		return w.elem.Wrap(c, e)
	}
	return data.Null, nil
}

func (w *pointerWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	e, err := w.elem.Unwrap(v)
	if err != nil {
		return _zero, err
	}
	return e.Addr(), nil
}
