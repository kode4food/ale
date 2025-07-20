package ffi

import (
	"encoding"
	"errors"
	"reflect"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

type marshaledWrapper struct {
	typ reflect.Type
}

var (
	textMarshaler   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	textUnmarshaler = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

func isMarshaledArray(t reflect.Type) bool {
	p := reflect.PointerTo(t)
	return p.Implements(textMarshaler) && p.Implements(textUnmarshaler)
}

func wrapMarshaled(t reflect.Type) (Wrapper, error) {
	return &marshaledWrapper{
		typ: t,
	}, nil
}

func (w *marshaledWrapper) Wrap(_ *Context, v reflect.Value) (ale.Value, error) {
	m := v.Interface().(encoding.TextMarshaler)
	s, err := m.MarshalText()
	if err != nil {
		return nil, err
	}
	return data.String(s), nil
}

func (w *marshaledWrapper) Unwrap(v ale.Value) (reflect.Value, error) {
	s, ok := v.(data.String)
	if !ok {
		return _zero, errors.New(ErrValueMustBeString)
	}
	out := reflect.New(w.typ)
	m := out.Interface().(encoding.TextUnmarshaler)
	err := m.UnmarshalText([]byte(s))
	if err != nil {
		return _zero, err
	}
	return out.Elem(), nil
}
