package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale"
)

type valueWrapper struct{}

// ErrMustImplementValue is raised when a value Unwrap call can't treat its
// source as a data.Value
var ErrMustImplementValue = errors.New("must implement value")

var dataValue = reflect.TypeOf((*ale.Value)(nil)).Elem()

func wrapDataValue(_ reflect.Type) (Wrapper, error) {
	return valueWrapper{}, nil
}

func (d valueWrapper) Wrap(_ *Context, v reflect.Value) (ale.Value, error) {
	if v, ok := v.Interface().(ale.Value); ok {
		return v, nil
	}
	return nil, ErrMustImplementValue
}

func (d valueWrapper) Unwrap(v ale.Value) (reflect.Value, error) {
	return reflect.ValueOf(v), nil
}
