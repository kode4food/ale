package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	floatWrapper  reflect.Kind
	intWrapper    reflect.Kind
	stringWrapper reflect.Kind
	boolWrapper   bool
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

func (f floatWrapper) Wrap(_ *WrapContext, v reflect.Value) data.Value {
	if !v.IsValid() {
		return data.Nil
	}
	return data.Float(v.Float())
}

func (f floatWrapper) Unwrap(_ *UnwrapContext, v data.Value) reflect.Value {
	switch reflect.Kind(f) {
	case reflect.Float32:
		if v == nil {
			return float32zero
		}
		return reflect.ValueOf(float32(v.(data.Float)))
	case reflect.Float64:
		if v == nil {
			return float64zero
		}
		return reflect.ValueOf(float64(v.(data.Float)))
	}
	panic("float kind is incorrect")
}

func (i intWrapper) Wrap(_ *WrapContext, v reflect.Value) data.Value {
	if !v.IsValid() {
		return data.Nil
	}
	return data.Integer(v.Int())
}

func (i intWrapper) Unwrap(_ *UnwrapContext, v data.Value) reflect.Value {
	switch reflect.Kind(i) {
	case reflect.Int64:
		if v == nil {
			return int64zero
		}
		return reflect.ValueOf(int64(v.(data.Integer)))
	case reflect.Int32:
		if v == nil {
			return int32zero
		}
		return reflect.ValueOf(int32(v.(data.Integer)))
	case reflect.Int16:
		if v == nil {
			return int16zero
		}
		return reflect.ValueOf(int16(v.(data.Integer)))
	case reflect.Int8:
		if v == nil {
			return int8zero
		}
		return reflect.ValueOf(int8(v.(data.Integer)))
	case reflect.Int:
		if v == nil {
			return intZero
		}
		return reflect.ValueOf(int(v.(data.Integer)))
	case reflect.Uint64:
		if v == nil {
			return uint64zero
		}
		return reflect.ValueOf(uint64(v.(data.Integer)))
	case reflect.Uint32:
		if v == nil {
			return uint32zero
		}
		return reflect.ValueOf(uint32(v.(data.Integer)))
	case reflect.Uint16:
		if v == nil {
			return uint16zero
		}
		return reflect.ValueOf(uint16(v.(data.Integer)))
	case reflect.Uint8:
		if v == nil {
			return uint8zero
		}
		return reflect.ValueOf(uint8(v.(data.Integer)))
	case reflect.Uint:
		if v == nil {
			return uintZero
		}
		return reflect.ValueOf(uint(v.(data.Integer)))
	}
	panic("int kind is incorrect")
}

func (stringWrapper) Wrap(_ *WrapContext, v reflect.Value) data.Value {
	if !v.IsValid() {
		return data.Nil
	}
	return data.String(v.Interface().(string))
}

func (stringWrapper) Unwrap(_ *UnwrapContext, v data.Value) reflect.Value {
	if v == nil {
		v = data.Nil
	}
	return reflect.ValueOf(v.String())
}

func (b boolWrapper) Wrap(_ *WrapContext, v reflect.Value) data.Value {
	if !v.IsValid() {
		return data.False
	}
	return data.Bool(v.Bool())
}

func (b boolWrapper) Unwrap(_ *UnwrapContext, v data.Value) reflect.Value {
	if v == nil {
		return boolZero
	}
	return reflect.ValueOf(bool(v.(data.Bool)))
}
