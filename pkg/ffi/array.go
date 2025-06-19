package ffi

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
)

type (
	arrayWrapper struct {
		typ  reflect.Type
		elem Wrapper
		len  int
	}

	byteArrayWrapper struct {
		typ reflect.Type
		len int
	}
)

const (
	// ErrValueMustBeSequence is raised when an Array Unwrap call can't treat its
	// source as a data.Sequence
	ErrValueMustBeSequence = "value must be a sequence"

	// ErrBadSliceLength is raised when an Array Unwrap call receives a slice
	// with a length that does not match the expected length of the array
	ErrBadSliceLength = "bad slice length: expected %d, got %d"
)

func makeWrappedArray(t reflect.Type) (Wrapper, error) {
	if t.Elem().Kind() == reflect.Uint8 {
		if isMarshaledArray(t) {
			return wrapMarshaled(t)
		}
		return wrapByteArray(t)
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
	out := make(data.Vector, w.len)
	for i := range w.len {
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
		return _zero, errors.New(ErrValueMustBeSequence)
	}
	in := sequence.ToVector(s)
	out := reflect.New(w.typ).Elem()
	for i, e := range in {
		v, err := w.elem.Unwrap(e)
		if err != nil {
			return _zero, err
		}
		out.Index(i).Set(v)
	}
	return out, nil
}

func wrapByteArray(t reflect.Type) (Wrapper, error) {
	return &byteArrayWrapper{
		typ: t,
		len: t.Len(),
	}, nil
}

func (b *byteArrayWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	if !v.CanAddr() {
		return b.slowWrap(v)
	}
	out := make(data.Bytes, b.len)
	copy(out, v.Bytes())
	return out, nil
}

func (b *byteArrayWrapper) slowWrap(v reflect.Value) (data.Value, error) {
	out := make(data.Bytes, b.len)
	for i := range b.len {
		out[i] = v.Index(i).Interface().(byte)
	}
	return out, nil
}

func (b *byteArrayWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	out := reflect.New(b.typ).Elem()
	if !out.CanAddr() {
		return b.slowUnwrap(v, out)
	}
	in, err := asByteArray(v)
	if err != nil {
		return _zero, err
	}
	copy(out.Bytes(), in)
	return out, nil
}

func (b *byteArrayWrapper) slowUnwrap(v data.Value, out reflect.Value) (reflect.Value, error) {
	in, err := asByteArray(v)
	if err != nil {
		return _zero, err
	}
	if len(in) != b.len {
		return _zero, fmt.Errorf(ErrBadSliceLength, b.len, len(in))
	}
	for i := range b.len {
		out.Index(i).Set(reflect.ValueOf(in[i]))
	}
	return out, nil
}

func asByteArray(v data.Value) ([]byte, error) {
	switch v := v.(type) {
	case data.Bytes:
		return v, nil
	case data.String:
		return []byte(v), nil
	default:
		return nil, errors.New(ErrValueMustBeString)
	}
}
