package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/data"
)

type (
	intWrapper   reflect.Kind
	int64Wrapper reflect.Kind
	int32Wrapper reflect.Kind
	int16Wrapper reflect.Kind
	int8Wrapper  reflect.Kind

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
		return intWrapper(k)
	case reflect.Int64:
		return int64Wrapper(k)
	case reflect.Int32:
		return int32Wrapper(k)
	case reflect.Int16:
		return int16Wrapper(k)
	case reflect.Int8:
		return int8Wrapper(k)
	default:
		panic(debug.ProgrammerError("int kind is incorrect"))
	}
}

func (intWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (intWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapInt[int](v)
}

func (int64Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (int64Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapInt[int64](v)
}

func (int32Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (int32Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapInt[int32](v)
}

func (int16Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (int16Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapInt[int16](v)
}

func (int8Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (int8Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapInt[int8](v)
}

func unwrapInt[T wrappableInts](v data.Value) (reflect.Value, error) {
	if v, ok := v.(data.Integer); ok {
		return reflect.ValueOf(T(v)), nil
	}
	return zero[T](), errors.New(ErrValueMustBeInteger)
}
