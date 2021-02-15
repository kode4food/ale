package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	stringWrapper reflect.Kind
	boolWrapper   bool
)

// Error messages
const (
	ErrValueMustBeBool   = "value must be a bool"
	ErrValueMustBeString = "value must be a string"
)

var (
	_stringWrapper stringWrapper
	_boolWrapper   boolWrapper

	stringZero = reflect.ValueOf("")
	boolTrue   = reflect.ValueOf(true)
	boolFalse  = reflect.ValueOf(false)
)

func makeWrappedBool(_ reflect.Type) (Wrapper, error) {
	return _boolWrapper, nil
}

func makeWrappedString(_ reflect.Type) (Wrapper, error) {
	return _stringWrapper, nil
}

func (stringWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.String(v.Interface().(string)), nil
}

func (stringWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if s, ok := v.(data.String); ok {
		return reflect.ValueOf(string(s)), nil
	}
	return stringZero, errors.New(ErrValueMustBeString)
}

func (w boolWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Bool(v.Bool()), nil
}

func (w boolWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if b, ok := v.(data.Bool); ok {
		if b {
			return boolTrue, nil
		}
		return boolFalse, nil
	}
	return boolFalse, errors.New(ErrValueMustBeBool)
}
