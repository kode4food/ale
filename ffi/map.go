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

func makeWrappedMap(t reflect.Type) (Wrapper, error) {
	if kw, err := wrapType(t.Key()); err != nil {
		return nil, err
	} else if vw, err := wrapType(t.Elem()); err != nil {
		return nil, err
	} else {
		return &mapWrapper{
			typ:   t,
			key:   kw,
			value: vw,
		}, nil
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
	out := make(data.Pairs, 0, v.Len())
	for pairs := v.MapRange(); pairs.Next(); {
		if k, err := m.key.Wrap(c, pairs.Key()); err != nil {
			return nil, err
		} else if v, err := m.value.Wrap(c, pairs.Value()); err != nil {
			return nil, err
		} else {
			out = append(out, data.NewCons(k, v))
		}
	}
	return data.NewObject(out...), nil
}

func (m *mapWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	in, err := sequence.ToObject(v.(data.Sequence))
	if err != nil {
		return _emptyValue, err
	}
	out := reflect.MakeMapWithSize(m.typ, in.Count())
	for f, r, ok := in.Split(); ok; f, r, ok = r.Split() {
		p := f.(data.Pair)
		k := p.Car()
		v := p.Cdr()
		if k, err := m.key.Unwrap(k); err != nil {
			return _emptyValue, err
		} else if v, err := m.value.Unwrap(v); err != nil {
			return _emptyValue, err
		} else {
			out.SetMapIndex(k, v)
		}
	}
	return out, nil
}
