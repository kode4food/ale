package ffi

import (
	"reflect"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

type (
	sliceWrapper struct {
		typ  reflect.Type
		elem Wrapper
	}

	byteSliceWrapper struct {
		typ reflect.Type
	}
)

func makeWrappedSlice(t reflect.Type) (Wrapper, error) {
	if t.Elem().Kind() == reflect.Uint8 {
		if isMarshaledArray(t) {
			return wrapMarshaled(t)
		}
		return wrapByteSlice(t)
	}
	w, err := WrapType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &sliceWrapper{
		typ:  t,
		elem: w,
	}, nil
}

func wrapByteSlice(t reflect.Type) (Wrapper, error) {
	return &byteSliceWrapper{
		typ: t,
	}, nil
}

func (w *sliceWrapper) Wrap(c *Context, v reflect.Value) (ale.Value, error) {
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

func (w *sliceWrapper) Unwrap(v ale.Value) (reflect.Value, error) {
	switch in := v.(type) {
	case data.Vector:
		return w.unwrapVector(in)
	case data.Counted:
		return w.unwrapCounted(in)
	case data.Sequence:
		return w.unwrapUncounted(in)
	default:
		return _zero, ErrValueMustBeSequence
	}
}

func (w *sliceWrapper) unwrapVector(in data.Vector) (reflect.Value, error) {
	inLen := len(in)
	out := reflect.MakeSlice(w.typ, inLen, inLen)
	for i, e := range in {
		v, err := w.elem.Unwrap(e)
		if err != nil {
			return _zero, err
		}
		out.Index(i).Set(v)
	}
	return out, nil
}

func (w *sliceWrapper) unwrapCounted(
	in data.Counted,
) (reflect.Value, error) {
	inLen := in.Count()
	out := reflect.MakeSlice(w.typ, inLen, inLen)
	var r data.Sequence = in
	for i := range inLen {
		var f ale.Value
		f, r, _ = r.Split()
		v, err := w.elem.Unwrap(f)
		if err != nil {
			return _zero, err
		}
		out.Index(i).Set(v)
	}
	return out, nil
}

func (w *sliceWrapper) unwrapUncounted(
	in data.Sequence,
) (reflect.Value, error) {
	out := reflect.MakeSlice(w.typ, 0, 0)
	for f, r, ok := in.Split(); ok; f, r, ok = r.Split() {
		v, err := w.elem.Unwrap(f)
		if err != nil {
			return _zero, err
		}
		out = reflect.Append(out, v)
	}
	return out, nil
}

func (w *byteSliceWrapper) Wrap(
	_ *Context, v reflect.Value,
) (ale.Value, error) {
	m := v.Interface().([]byte)
	return data.Bytes(m), nil
}

func (w *byteSliceWrapper) Unwrap(v ale.Value) (reflect.Value, error) {
	return asValueOf[[]byte](v)
}
