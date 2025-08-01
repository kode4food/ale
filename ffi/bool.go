package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

type boolWrapper struct{}

// ErrValueMustBeBool is raised when a boolean Unwrap call can't treat its
// source as a data.Bool
const ErrValueMustBeBool = "value must be a bool"

var (
	boolTrue  = reflect.ValueOf(true)
	boolFalse = reflect.ValueOf(false)
)

func (boolWrapper) Wrap(_ *Context, v reflect.Value) (ale.Value, error) {
	return data.Bool(v.Bool()), nil
}

func (boolWrapper) Unwrap(v ale.Value) (reflect.Value, error) {
	if b, ok := v.(data.Bool); ok {
		if b {
			return boolTrue, nil
		}
		return boolFalse, nil
	}
	return _zero, errors.New(ErrValueMustBeBool)
}
