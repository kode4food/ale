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

func (f floatWrapper) Wrap(v reflect.Value) data.Value {
	return data.Float(v.Float())
}

func (f floatWrapper) Unwrap(v data.Value) reflect.Value {
	switch reflect.Kind(f) {
	case reflect.Float32:
		return reflect.ValueOf(float32(v.(data.Float)))
	case reflect.Float64:
		return reflect.ValueOf(float64(v.(data.Float)))
	}
	panic("float kind is incorrect")
}

func (i intWrapper) Wrap(v reflect.Value) data.Value {
	return data.Integer(v.Int())
}

func (i intWrapper) Unwrap(v data.Value) reflect.Value {
	switch reflect.Kind(i) {
	case reflect.Int64:
		return reflect.ValueOf(int64(v.(data.Integer)))
	case reflect.Int32:
		return reflect.ValueOf(int32(v.(data.Integer)))
	case reflect.Int16:
		return reflect.ValueOf(int16(v.(data.Integer)))
	case reflect.Int8:
		return reflect.ValueOf(int8(v.(data.Integer)))
	case reflect.Int:
		return reflect.ValueOf(int(v.(data.Integer)))
	case reflect.Uint64:
		return reflect.ValueOf(uint64(v.(data.Integer)))
	case reflect.Uint32:
		return reflect.ValueOf(uint32(v.(data.Integer)))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(v.(data.Integer)))
	case reflect.Uint8:
		return reflect.ValueOf(uint8(v.(data.Integer)))
	case reflect.Uint:
		return reflect.ValueOf(uint(v.(data.Integer)))
	}
	panic("int kind is incorrect")
}

func (stringWrapper) Wrap(v reflect.Value) data.Value {
	return data.String(v.Interface().(string))
}

func (stringWrapper) Unwrap(v data.Value) reflect.Value {
	return reflect.ValueOf(v.String())
}

func (b boolWrapper) Wrap(v reflect.Value) data.Value {
	return data.Bool(v.Bool())
}

func (b boolWrapper) Unwrap(v data.Value) reflect.Value {
	return reflect.ValueOf(bool(v.(data.Bool)))
}
