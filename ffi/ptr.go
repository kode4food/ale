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

func (p *pointerWrapper) Wrap(c *WrapContext, v reflect.Value) data.Value {
	if !v.IsValid() {
		return data.Nil
	}
	return p.elem.Wrap(c, v.Elem())
}

func (p *pointerWrapper) Unwrap(c *UnwrapContext, v data.Value) reflect.Value {
	return p.elem.Unwrap(c, v).Addr()
}
