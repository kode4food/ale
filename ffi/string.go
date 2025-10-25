package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

type stringWrapper struct{}

// ErrValueMustBeString is raised when a string Unwrap call can't treat its
// source as a data.String
var ErrValueMustBeString = errors.New("value must be a byte slice or string")

func (stringWrapper) Wrap(_ *Context, v reflect.Value) (ale.Value, error) {
	return data.String(v.String()), nil
}

func (stringWrapper) Unwrap(v ale.Value) (reflect.Value, error) {
	return asValueOf[string](v)
}

func asValueOf[T string | []byte](v ale.Value) (reflect.Value, error) {
	switch v := v.(type) {
	case data.String:
		return reflect.ValueOf(T(v)), nil
	case data.Bytes:
		return reflect.ValueOf(T(v)), nil
	default:
		return _zero, ErrValueMustBeString
	}
}
