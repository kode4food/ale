package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

type mapWrapper struct {
	typ   reflect.Type
	key   Wrapper
	value Wrapper
}

func makeWrappedMap(t reflect.Type) Wrapper {
	return &mapWrapper{
		typ:   t,
		key:   wrapType(t.Key()),
		value: wrapType(t.Elem()),
	}
}

func (m *mapWrapper) Wrap(c *WrapContext, v reflect.Value) data.Value {
	if r, ok := c.Get(v); ok {
		return r
	}
	out := make(data.Object, v.Len())
	c.Put(v, out)
	pairs := v.MapRange()
	for pairs.Next() {
		k := pairs.Key()
		v := pairs.Value()
		out[m.key.Wrap(c, k)] = m.value.Wrap(c, v)
	}
	return out
}

func (m *mapWrapper) Unwrap(c *UnwrapContext, v data.Value) reflect.Value {
	if r, ok := c.Get(v); ok {
		return r
	}
	in := sequence.ToObject(v.(data.Sequence))
	out := reflect.MakeMapWithSize(m.typ, len(in))
	c.Put(v, out)
	for k, v := range in {
		out.SetMapIndex(m.key.Unwrap(c, k), m.value.Unwrap(c, v))
	}
	return out
}
