package ffi

import (
	"errors"
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
	kw, err := wrapType(t.Key())
	if err != nil {
		return nil, err
	}
	vw, err := wrapType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &mapWrapper{
		typ:   t,
		key:   kw,
		value: vw,
	}, nil
}

func (w *mapWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	if !v.IsValid() {
		return data.Nil, nil
	}
	c, err := c.Push(v)
	if err != nil {
		return data.Nil, err
	}
	out := make(data.Pairs, 0, v.Len())
	for pairs := v.MapRange(); pairs.Next(); {
		k, err := w.key.Wrap(c, pairs.Key())
		if err != nil {
			return data.Nil, err
		}
		v, err := w.value.Wrap(c, pairs.Value())
		if err != nil {
			return data.Nil, err
		}
		out = append(out, data.NewCons(k, v))
	}
	return data.NewObject(out...), nil
}

func (w *mapWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if s, ok := v.(data.Sequence); ok {
		in, err := sequence.ToObject(s)
		if err != nil {
			return _emptyValue, err
		}
		out := reflect.MakeMapWithSize(w.typ, in.Count())
		for f, r, ok := in.Split(); ok; f, r, ok = r.Split() {
			p := f.(data.Pair)
			k := p.Car()
			v := p.Cdr()
			uk, err := w.key.Unwrap(k)
			if err != nil {
				return _emptyValue, err
			}
			uv, err := w.value.Unwrap(v)
			if err != nil {
				return _emptyValue, err
			}
			out.SetMapIndex(uk, uv)
		}
		return out, nil
	}
	return _emptyValue, errors.New(ErrValueMustBeSequence)
}
