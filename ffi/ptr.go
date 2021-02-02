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
	c, err := c.Push(v)
	if err != nil {
		return nil, err
	}
	e := v.Elem()
	if e.IsValid() {
		return p.elem.Wrap(c, e)
	}
	return data.Nil, nil
}

func (p *pointerWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	e, err := p.elem.Unwrap(v)
	if err != nil {
		return _emptyValue, nil
	}
	return e.Addr(), nil
}
