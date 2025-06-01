package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/data"
)

type (
	float32Wrapper reflect.Kind
	float64Wrapper reflect.Kind
)

// ErrValueMustBeFloat is raised when a float Unwrap call can't treat its
// source as a data.Integer or data.Float
const ErrValueMustBeFloat = "value must be a float"

func makeWrappedFloat(t reflect.Type) Wrapper {
	switch k := t.Kind(); k {
	case reflect.Float32:
		return float32Wrapper(k)
	case reflect.Float64:
		return float64Wrapper(k)
	default:
		panic(debug.ProgrammerError("float kind is incorrect"))
	}
}

func (float32Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Float(v.Float()), nil
}

func (float32Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapFloat[float32](v)
}

func (float64Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Float(v.Float()), nil
}

func (float64Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapFloat[float64](v)
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

func unwrapFloat[T ~float32 | ~float64](v data.Value) (reflect.Value, error) {
	if f, ok := makeFloat64(v); ok {
		return reflect.ValueOf(T(f)), nil
	}
	return reflect.Value{}, errors.New(ErrValueMustBeFloat)
}
