package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	float32Wrapper reflect.Kind
	float64Wrapper reflect.Kind
)

// Error messages
const (
	ErrValueMustBeFloat = "value must be a float"

	errIncorrectFloatKind = "float kind is incorrect"
)

var (
	float32zero = reflect.ValueOf(float32(0))
	float64zero = reflect.ValueOf(float64(0))
)

func makeWrappedFloat(t reflect.Type) (Wrapper, error) {
	k := t.Kind()
	switch k {
	case reflect.Float32:
		return float32Wrapper(k), nil
	case reflect.Float64:
		return float64Wrapper(k), nil
	default:
		// Programmer error
		panic(errors.New(errIncorrectFloatKind))
	}
}

func (float32Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Float(v.Float()), nil
}

func (float32Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if f, ok := v.(data.Float); ok {
		return reflect.ValueOf(float32(f)), nil
	}
	return float32zero, errors.New(ErrValueMustBeFloat)
}

func (float64Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Float(v.Float()), nil
}

func (float64Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if f, ok := v.(data.Float); ok {
		return reflect.ValueOf(float64(f)), nil
	}
	return float64zero, errors.New(ErrValueMustBeFloat)
}
