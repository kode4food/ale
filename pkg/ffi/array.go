package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
)

type arrayWrapper struct {
	typ  reflect.Type
	elem Wrapper
	len  int
}

// ErrValueMustBeSequence is raised when an Array Unwrap call can't treat its
// source as a data.Sequence
const ErrValueMustBeSequence = "value must be a sequence"

func makeWrappedArray(t reflect.Type) (Wrapper, error) {
	if isMarshaledByteArray(t) {
		return wrapMarshaledByteArray(t)
	}
	w, err := WrapType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &arrayWrapper{
		typ:  t,
		len:  t.Len(),
		elem: w,
	}, nil
}

func (w *arrayWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	vLen := v.Len()
	out := make(data.Vector, vLen)
	for i := range vLen {
		elem, err := w.elem.Wrap(c, v.Index(i))
		if err != nil {
			return data.Null, err
		}
		out[i] = elem
	}
	return out, nil
}

func (w *arrayWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	s, ok := v.(data.Sequence)
	if !ok {
		return _emptyValue, errors.New(ErrValueMustBeSequence)
	}
	in := sequence.ToVector(s)
	out := reflect.New(w.typ).Elem()
	for i, e := range in {
		v, err := w.elem.Unwrap(e)
		if err != nil {
			return _emptyValue, err
		}
		out.Index(i).Set(v)
	}
	return out, nil
}
