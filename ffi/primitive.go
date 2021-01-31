package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	floatWrapper  reflect.Kind
	intWrapper    reflect.Kind
	stringWrapper reflect.Kind
	boolWrapper   bool
)

// Error messages
const (
	ErrIncorrectFloatKind = "float kind is incorrect"
	ErrIncorrectIntKind   = "int kind is incorrect"
)

var (
	_stringWrapper stringWrapper
	_boolWrapper   boolWrapper

	float32zero = reflect.ValueOf(float32(0))
	float64zero = reflect.ValueOf(float64(0))

	int64zero = reflect.ValueOf(int64(0))
	int32zero = reflect.ValueOf(int32(0))
	int16zero = reflect.ValueOf(int16(0))
	int8zero  = reflect.ValueOf(int8(0))
	intZero   = reflect.ValueOf(int(0))

	uint64zero = reflect.ValueOf(uint64(0))
	uint32zero = reflect.ValueOf(uint32(0))
	uint16zero = reflect.ValueOf(uint16(0))
	uint8zero  = reflect.ValueOf(uint8(0))
	uintZero   = reflect.ValueOf(uint(0))

	boolZero = reflect.ValueOf(false)
)

func makeWrappedFloat(t reflect.Type) Wrapper {
	return floatWrapper(t.Kind())
}

func makeWrappedInt(t reflect.Type) Wrapper {
	return intWrapper(t.Kind())
}

func makeWrappedBool(_ reflect.Type) Wrapper {
	return _boolWrapper
}

func makeWrappedString(_ reflect.Type) Wrapper {
	return _stringWrapper
}

func (f floatWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	if !v.IsValid() {
		return data.Nil, nil
	}
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
	return emptyReflectValue, errors.New(ErrIncorrectFloatKind)
}

func (i intWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	if !v.IsValid() {
		return data.Nil, nil
	}
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
	return emptyReflectValue, errors.New(ErrIncorrectIntKind)
}

func (stringWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	if !v.IsValid() {
		return data.Nil, nil
	}
	return data.String(v.Interface().(string)), nil
}

func (stringWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if v == nil {
		v = data.Nil
	}
	return reflect.ValueOf(v.String()), nil
}

func (b boolWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	if !v.IsValid() {
		return data.False, nil
	}
	return data.Bool(v.Bool()), nil
}

func (b boolWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if v == nil {
		return boolZero, nil
	}
	return reflect.ValueOf(bool(v.(data.Bool))), nil
}
