package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type boolWrapper reflect.Kind

// Error messages
const (
	ErrValueMustBeBool = "value must be a bool"
)

var (
	boolTrue  = reflect.ValueOf(true)
	boolFalse = reflect.ValueOf(false)
)

func makeWrappedBool(t reflect.Type) (Wrapper, error) {
	return boolWrapper(t.Kind()), nil
}

func (boolWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Bool(v.Bool()), nil
}

func (boolWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if b, ok := v.(data.Bool); ok {
		if b {
			return boolTrue, nil
		}
		return boolFalse, nil
	}
	return boolFalse, errors.New(ErrValueMustBeBool)
}
