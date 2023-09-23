package ffi

import (
	"fmt"
	"reflect"

	"github.com/kode4food/ale/data"
)

type valueWrapper struct{}

// Error messages
const (
	ErrMustImplementValue = "must implement value"
)

var dataValue = reflect.TypeOf((*data.Value)(nil)).Elem()

func wrapDataValue(_ reflect.Type) (Wrapper, error) {
	return valueWrapper{}, nil
}

func (d valueWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	if v, ok := v.Interface().(data.Value); ok {
		return v, nil
	}
	return nil, fmt.Errorf(ErrMustImplementValue)
}

func (d valueWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return reflect.ValueOf(v), nil
}
