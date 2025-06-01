package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/data"
)

type (
	intWrapper[T wrappableInts] struct{}

	wrappableInts interface {
		~int | ~uint8 | ~int8 | ~uint16 | ~int16 | ~uint32 | ~int32 | ~int64
	}
)

// ErrValueMustBeInteger is raised when an integer Unwrap call can't treat its
// source as a data.Integer
const ErrValueMustBeInteger = "value must be an integer"

func makeWrappedInt(t reflect.Type) Wrapper {
	switch k := t.Kind(); k {
	case reflect.Int:
		return intWrapper[int]{}
	case reflect.Int64:
		return intWrapper[int64]{}
	case reflect.Int32:
		return intWrapper[int32]{}
	case reflect.Int16:
		return intWrapper[int16]{}
	case reflect.Int8:
		return intWrapper[int8]{}
	default:
		panic(debug.ProgrammerError("int kind is incorrect"))
	}
}

func (intWrapper[_]) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (intWrapper[T]) Unwrap(v data.Value) (reflect.Value, error) {
	if v, ok := v.(data.Integer); ok {
		return reflect.ValueOf(T(v)), nil
	}
	return zero[T](), errors.New(ErrValueMustBeInteger)
}
