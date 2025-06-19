package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/pkg/data"
)

type stringWrapper struct{}

// ErrValueMustBeString is raised when a string Unwrap call can't treat its
// source as a data.String
const ErrValueMustBeString = "value must be a byte slice or string"

func (stringWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.String(v.String()), nil
}

func (stringWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return asValueOf[string](v)
}

func asValueOf[T string | []byte](v data.Value) (reflect.Value, error) {
	switch v := v.(type) {
	case data.String:
		return reflect.ValueOf(T(v)), nil
	case data.Bytes:
		return reflect.ValueOf(T(v)), nil
	default:
		return _zero, errors.New(ErrValueMustBeString)
	}
}
