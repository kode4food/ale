package ffi

import (
	"reflect"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

type mapWrapper struct {
	typ   reflect.Type
	key   Wrapper
	value Wrapper
}

func makeWrappedMap(t reflect.Type) (Wrapper, error) {
	kw, err := WrapType(t.Key())
	if err != nil {
		return nil, err
	}
	vw, err := WrapType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &mapWrapper{
		typ:   t,
		key:   kw,
		value: vw,
	}, nil
}

func (w *mapWrapper) Wrap(c *Context, v reflect.Value) (ale.Value, error) {
	if !v.IsValid() {
		return data.Null, nil
	}
	c, err := c.Push(v)
	if err != nil {
		return data.Null, err
	}
	out := make(data.Pairs, 0, v.Len())
	for pairs := v.MapRange(); pairs.Next(); {
		k, err := w.key.Wrap(c, pairs.Key())
		if err != nil {
			return data.Null, err
		}
		v, err := w.value.Wrap(c, pairs.Value())
		if err != nil {
			return data.Null, err
		}
		out = append(out, data.NewCons(k, v))
	}
	return data.NewObject(out...), nil
}

func (w *mapWrapper) Unwrap(v ale.Value) (reflect.Value, error) {
	s, ok := v.(data.Sequence)
	if !ok {
		return _zero, ErrValueMustBeSequence
	}
	in, err := sequence.ToObject(s)
	if err != nil {
		return _zero, err
	}
	out := reflect.MakeMapWithSize(w.typ, in.Count())
	for f, r, ok := in.Split(); ok; f, r, ok = r.Split() {
		p := f.(data.Pair)
		k := p.Car()
		v := p.Cdr()
		uk, err := w.key.Unwrap(k)
		if err != nil {
			return _zero, err
		}
		uv, err := w.value.Unwrap(v)
		if err != nil {
			return _zero, err
		}
		out.SetMapIndex(uk, uv)
	}
	return out, nil
}
