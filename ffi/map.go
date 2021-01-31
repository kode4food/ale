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

func (m *mapWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	if !v.IsValid() {
		return data.Nil, nil
	}
	c, err := c.Push(v)
	if err != nil {
		return nil, err
	}
	out := make(data.Object, v.Len())
	pairs := v.MapRange()
	for pairs.Next() {
		k, err := m.key.Wrap(c, pairs.Key())
		if err != nil {
			return nil, err
		}
		v, err := m.value.Wrap(c, pairs.Value())
		if err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, nil
}

func (m *mapWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	in := sequence.ToObject(v.(data.Sequence))
	out := reflect.MakeMapWithSize(m.typ, len(in))
	for k, v := range in {
		k, err := m.key.Unwrap(k)
		if err != nil {
			return emptyReflectValue, err
		}
		v, err := m.value.Unwrap(v)
		if err != nil {
			return emptyReflectValue, err
		}
		out.SetMapIndex(k, v)
	}
	return out, nil
}
