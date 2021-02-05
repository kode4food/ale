package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	intWrapper  reflect.Kind
	uintWrapper reflect.Kind
)

// Error messages
const (
	errIncorrectIntKind         = "int kind is incorrect"
	errIncorrectUnsignedIntKind = "uint kind is incorrect"
)

var (
	int64zero = reflect.ValueOf(int64(0))
	int32zero = reflect.ValueOf(int32(0))
	int16zero = reflect.ValueOf(int16(0))
	int8zero  = reflect.ValueOf(int8(0))
	intZero   = reflect.ValueOf(0)

	uint64zero = reflect.ValueOf(uint64(0))
	uint32zero = reflect.ValueOf(uint32(0))
	uint16zero = reflect.ValueOf(uint16(0))
	uint8zero  = reflect.ValueOf(uint8(0))
	uintZero   = reflect.ValueOf(uint(0))
)

func makeWrappedInt(t reflect.Type) (Wrapper, error) {
	return intWrapper(t.Kind()), nil
}

func makeWrappedUnsignedInt(t reflect.Type) (Wrapper, error) {
	return uintWrapper(t.Kind()), nil
}

func (i intWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (i intWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	switch reflect.Kind(i) {
	case reflect.Int64:
		if v == nil {
			return int64zero, nil
		}
		return reflect.ValueOf(int64(v.(data.Integer))), nil
	case reflect.Int32:
		if v == nil {
			return int32zero, nil
		}
		return reflect.ValueOf(int32(v.(data.Integer))), nil
	case reflect.Int16:
		if v == nil {
			return int16zero, nil
		}
		return reflect.ValueOf(int16(v.(data.Integer))), nil
	case reflect.Int8:
		if v == nil {
			return int8zero, nil
		}
		return reflect.ValueOf(int8(v.(data.Integer))), nil
	case reflect.Int:
		if v == nil {
			return intZero, nil
		}
		return reflect.ValueOf(int(v.(data.Integer))), nil
	}
	// Programmer error
	panic(errors.New(errIncorrectIntKind))
}

func (i uintWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Uint()), nil
}

func (i uintWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	switch reflect.Kind(i) {
	case reflect.Uint64:
		if v == nil {
			return uint64zero, nil
		}
		return reflect.ValueOf(uint64(v.(data.Integer))), nil
	case reflect.Uint32:
		if v == nil {
			return uint32zero, nil
		}
		return reflect.ValueOf(uint32(v.(data.Integer))), nil
	case reflect.Uint16:
		if v == nil {
			return uint16zero, nil
		}
		return reflect.ValueOf(uint16(v.(data.Integer))), nil
	case reflect.Uint8:
		if v == nil {
			return uint8zero, nil
		}
		return reflect.ValueOf(uint8(v.(data.Integer))), nil
	case reflect.Uint:
		if v == nil {
			return uintZero, nil
		}
		return reflect.ValueOf(uint(v.(data.Integer))), nil
	}
	// Programmer error
	panic(errors.New(errIncorrectUnsignedIntKind))
}
