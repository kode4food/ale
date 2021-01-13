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

func (p *pointerWrapper) Wrap(v reflect.Value) data.Value {
	return p.elem.Wrap(v.Elem())
}

func (p *pointerWrapper) Unwrap(v data.Value) reflect.Value {
	return p.elem.Unwrap(v).Addr()
}
