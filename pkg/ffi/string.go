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

var stringZero = reflect.ValueOf("")

func (stringWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.String(v.String()), nil
}

func (stringWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	switch v := v.(type) {
	case data.String:
		return reflect.ValueOf(string(v)), nil
	case data.Bytes:
		return reflect.ValueOf(string(v)), nil
	default:
		return stringZero, errors.New(ErrValueMustBeString)
	}
}
