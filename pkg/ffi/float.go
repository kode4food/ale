package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/pkg/data"
)

type floatWrapper[T ~float32 | ~float64] struct{}

// ErrValueMustBeFloat is raised when a float Unwrap call can't treat its
// source as a data.Integer or data.Float
const ErrValueMustBeFloat = "value must be a float"

func (floatWrapper[_]) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Float(v.Float()), nil
}

func (floatWrapper[T]) Unwrap(v data.Value) (reflect.Value, error) {
	if f, ok := makeFloat64(v); ok {
		return reflect.ValueOf(T(f)), nil
	}
	return reflect.Value{}, errors.New(ErrValueMustBeFloat)
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
