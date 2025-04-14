package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/pkg/data"
)

type sliceWrapper struct {
	typ  reflect.Type
	elem Wrapper
}

func makeWrappedSlice(t reflect.Type) (Wrapper, error) {
	w, err := WrapType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &sliceWrapper{
		typ:  t,
		elem: w,
	}, nil
}

func (w *sliceWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	c, err := c.Push(v)
	if err != nil {
		return data.Null, err
	}
	vLen := v.Len()
	out := make(data.Vector, vLen)
	for i := range vLen {
		v, err := w.elem.Wrap(c, v.Index(i))
		if err != nil {
			return data.Null, err
		}
		out[i] = v
	}
	return out, nil
}

func (w *sliceWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	switch in := v.(type) {
	case data.Vector:
		return w.unwrapVector(in)
	case data.CountedSequence:
		return w.unwrapCounted(in)
	case data.Sequence:
		return w.unwrapUncounted(in)
	default:
		return _emptyValue, errors.New(ErrValueMustBeSequence)
	}
}

func (w *sliceWrapper) unwrapVector(in data.Vector) (reflect.Value, error) {
	inLen := len(in)
	out := reflect.MakeSlice(w.typ, inLen, inLen)
	for i, e := range in {
		v, err := w.elem.Unwrap(e)
		if err != nil {
			return _emptyValue, err
		}
		out.Index(i).Set(v)
	}
	return out, nil
}

func (w *sliceWrapper) unwrapCounted(in data.CountedSequence) (reflect.Value, error) {
	inLen := int(in.Count())
	out := reflect.MakeSlice(w.typ, inLen, inLen)
	var r data.Sequence = in
	for i := range inLen {
		var f data.Value
		f, r, _ = r.Split()
		v, err := w.elem.Unwrap(f)
		if err != nil {
			return _emptyValue, err
		}
		out.Index(i).Set(v)
	}
	return out, nil
}

func (w *sliceWrapper) unwrapUncounted(in data.Sequence) (reflect.Value, error) {
	out := reflect.MakeSlice(w.typ, 0, 0)
	for f, r, ok := in.Split(); ok; f, r, ok = r.Split() {
		v, err := w.elem.Unwrap(f)
		if err != nil {
			return _emptyValue, err
		}
		out = reflect.Append(out, v)
	}
	return out, nil
}
