package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	stringWrapper reflect.Kind
	boolWrapper   bool
)

var (
	_stringWrapper stringWrapper
	_boolWrapper   boolWrapper

	boolZero = reflect.ValueOf(false)
)

func makeWrappedBool(_ reflect.Type) Wrapper {
	return _boolWrapper
}

func makeWrappedString(_ reflect.Type) Wrapper {
	return _stringWrapper
}

func (stringWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.String(v.Interface().(string)), nil
}

func (stringWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if v == nil {
		v = data.Nil
	}
	return reflect.ValueOf(v.String()), nil
}

func (b boolWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Bool(v.Bool()), nil
}

func (b boolWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if v == nil {
		return boolZero, nil
	}
	return reflect.ValueOf(bool(v.(data.Bool))), nil
}
