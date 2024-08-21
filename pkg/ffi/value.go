package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/pkg/data"
)

type valueWrapper struct{}

// ErrMustImplementValue is raised when a value Unwrap call can't treat its
// source as a data.Value
const ErrMustImplementValue = "must implement value"

var dataValue = reflect.TypeOf((*data.Value)(nil)).Elem()

func wrapDataValue(_ reflect.Type) (Wrapper, error) {
	return valueWrapper{}, nil
}

func (d valueWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	if v, ok := v.Interface().(data.Value); ok {
		return v, nil
	}
	return nil, errors.New(ErrMustImplementValue)
}

func (d valueWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return reflect.ValueOf(v), nil
}
