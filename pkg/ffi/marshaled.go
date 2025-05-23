package ffi

import (
	"encoding"
	"errors"
	"reflect"

	"github.com/kode4food/ale/pkg/data"
)

type byteArrayWrapper struct {
	typ reflect.Type
}

var (
	textMarshaler   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	textUnmarshaler = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

func isMarshaledByteArray(t reflect.Type) bool {
	if t.Kind() != reflect.Array || t.Elem().Kind() != reflect.Uint8 {
		return false
	}
	p := reflect.PointerTo(t)
	return p.Implements(textUnmarshaler) && p.Implements(textMarshaler)
}

func wrapMarshaledByteArray(t reflect.Type) (Wrapper, error) {
	return &byteArrayWrapper{
		typ: t,
	}, nil
}

func (w *byteArrayWrapper) Wrap(
	_ *Context, v reflect.Value,
) (data.Value, error) {
	m := v.Interface().(encoding.TextMarshaler)
	s, err := m.MarshalText()
	if err != nil {
		return nil, err
	}
	return data.String(s), nil
}

func (w *byteArrayWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	s, ok := v.(data.String)
	if !ok {
		return _emptyValue, errors.New(ErrValueMustBeString)
	}
	out := reflect.New(w.typ)
	m := out.Interface().(encoding.TextUnmarshaler)
	err := m.UnmarshalText([]byte(s))
	if err != nil {
		return _emptyValue, err
	}
	return out.Elem(), nil
}
