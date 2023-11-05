package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type stringWrapper reflect.Kind

// ErrValueMustBeString is raised when an string Unwrap call can't treat its
// source as a data.String
const ErrValueMustBeString = "value must be a string"

var stringZero = reflect.ValueOf("")

func makeWrappedString(t reflect.Type) (Wrapper, error) {
	return stringWrapper(t.Kind()), nil
}

func (stringWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.String(v.String()), nil
}

func (stringWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if s, ok := v.(data.String); ok {
		return reflect.ValueOf(string(s)), nil
	}
	return stringZero, errors.New(ErrValueMustBeString)
}
