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

func (m *mapWrapper) Wrap(v reflect.Value) data.Value {
	out := make(data.Object, v.Len())
	pairs := v.MapRange()
	for pairs.Next() {
		k := pairs.Key()
		v := pairs.Value()
		out[m.key.Wrap(k)] = m.value.Wrap(v)
	}
	return out
}

func (m *mapWrapper) Unwrap(v data.Value) reflect.Value {
	in := sequence.ToObject(v.(data.Sequence))
	out := reflect.MakeMapWithSize(m.typ, len(in))
	for k, v := range in {
		out.SetMapIndex(m.key.Unwrap(k), m.value.Unwrap(v))
	}
	return out
}
