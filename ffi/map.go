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
		if k, err := w.key.Wrap(c, pairs.Key()); err != nil {
			return data.Nil, err
		} else if v, err := w.value.Wrap(c, pairs.Value()); err != nil {
			return data.Nil, err
		} else {
			out = append(out, data.NewCons(k, v))
		}
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
			if k, err := w.key.Unwrap(k); err != nil {
				return _emptyValue, err
			} else if v, err := w.value.Unwrap(v); err != nil {
				return _emptyValue, err
			} else {
				out.SetMapIndex(k, v)
			}
		}
		return out, nil
	}
	return _emptyValue, errors.New(ErrValueMustBeSequence)
}
