package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type stringWrapper reflect.Kind

// Error messages
const (
	ErrValueMustBeString = "value must be a string"
)

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
