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
)

// ErrValueMustBeInteger is raised when an integer Unwrap call can't treat its
// source as a data.Integer
const ErrValueMustBeInteger = "value must be an integer"

var (
	intZero   = reflect.ValueOf(0)
	int64zero = reflect.ValueOf(int64(0))
	int32zero = reflect.ValueOf(int32(0))
	int16zero = reflect.ValueOf(int16(0))
	int8zero  = reflect.ValueOf(int8(0))
)

func makeWrappedInt(t reflect.Type) (Wrapper, error) {
	switch k := t.Kind(); k {
	case reflect.Int:
		return intWrapper(k), nil
	case reflect.Int64:
		return int64Wrapper(k), nil
	case reflect.Int32:
		return int32Wrapper(k), nil
	case reflect.Int16:
		return int16Wrapper(k), nil
	case reflect.Int8:
		return int8Wrapper(k), nil
	default:
		panic(debug.ProgrammerError("int kind is incorrect"))
	}
}

func (intWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (intWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if i, ok := v.(data.Integer); ok {
		return reflect.ValueOf(int(i)), nil
	}
	return intZero, errors.New(ErrValueMustBeInteger)
}

func (int64Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (int64Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if i, ok := v.(data.Integer); ok {
		return reflect.ValueOf(int64(i)), nil
	}
	return int64zero, errors.New(ErrValueMustBeInteger)
}

func (int32Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (int32Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if i, ok := v.(data.Integer); ok {
		return reflect.ValueOf(int32(i)), nil
	}
	return int32zero, errors.New(ErrValueMustBeInteger)
}

func (int16Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (int16Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if i, ok := v.(data.Integer); ok {
		return reflect.ValueOf(int16(i)), nil
	}
	return int16zero, errors.New(ErrValueMustBeInteger)
}

func (int8Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (int8Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if i, ok := v.(data.Integer); ok {
		return reflect.ValueOf(int8(i)), nil
	}
	return int8zero, errors.New(ErrValueMustBeInteger)
}
