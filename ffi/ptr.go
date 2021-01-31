package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
)

type pointerWrapper struct {
	elem Wrapper
}

func makeWrappedPointer(t reflect.Type) Wrapper {
	return &pointerWrapper{
		elem: wrapType(t.Elem()),
	}
}

func (p *pointerWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	if !v.IsValid() {
		return data.Nil, nil
	}
	c, err := c.Push(v)
	if err != nil {
		return nil, err
	}
	return p.elem.Wrap(c, v.Elem())
}

func (p *pointerWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	e, err := p.elem.Unwrap(v)
	if err != nil {
		return _emptyValue, nil
	}
	return e.Addr(), nil
}
