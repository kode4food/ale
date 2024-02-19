package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/debug"
)

type (
	float32Wrapper reflect.Kind
	float64Wrapper reflect.Kind
)

// ErrValueMustBeFloat is raised when a float Unwrap call can't treat its
// source as a data.Integer or data.Float
const ErrValueMustBeFloat = "value must be a float"

var (
	float32zero = reflect.ValueOf(float32(0))
	float64zero = reflect.ValueOf(float64(0))
)

func makeWrappedFloat(t reflect.Type) (Wrapper, error) {
	switch k := t.Kind(); k {
	case reflect.Float32:
		return float32Wrapper(k), nil
	case reflect.Float64:
		return float64Wrapper(k), nil
	default:
		panic(debug.ProgrammerError("float kind is incorrect"))
	}
}

func (float32Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Float(v.Float()), nil
}

func (float32Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if f, ok := makeFloat64(v); ok {
		return reflect.ValueOf(float32(f)), nil
	}
	return float32zero, errors.New(ErrValueMustBeFloat)
}

func (float64Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Float(v.Float()), nil
}

func (float64Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if f, ok := makeFloat64(v); ok {
		return reflect.ValueOf(f), nil
	}
	return float64zero, errors.New(ErrValueMustBeFloat)
}

func makeFloat64(v data.Value) (float64, bool) {
	switch v := v.(type) {
	case data.Integer:
		return float64(v), true
	case data.Float:
		return float64(v), true
	default:
		return 0, false
	}
}
