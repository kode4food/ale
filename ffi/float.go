package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type floatWrapper reflect.Kind

// Error messages
const (
	errIncorrectFloatKind = "float kind is incorrect"
)

var (
	float32zero = reflect.ValueOf(float32(0))
	float64zero = reflect.ValueOf(float64(0))
)

func makeWrappedFloat(t reflect.Type) Wrapper {
	return floatWrapper(t.Kind())
}

func (f floatWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Float(v.Float()), nil
}

func (f floatWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	switch reflect.Kind(f) {
	case reflect.Float32:
		if v == nil {
			return float32zero, nil
		}
		return reflect.ValueOf(float32(v.(data.Float))), nil
	case reflect.Float64:
		if v == nil {
			return float64zero, nil
		}
		return reflect.ValueOf(float64(v.(data.Float))), nil
	}
	panic(errors.New(errIncorrectFloatKind))
}
