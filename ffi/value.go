package ffi

import (
	"fmt"
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	dataValueWrapper bool
	valueWrapper     bool

	wrappedValue struct {
		elem reflect.Value
	}
)

const wrappedValueType = data.Name("wrapped-value")

var (
	dataValueType     = reflect.TypeOf((*data.Value)(nil)).Elem()
	_dataValueWrapper dataValueWrapper
	_valueWrapper     valueWrapper
)

func makeWrappedValue(t reflect.Type) Wrapper {
	if t.Implements(dataValueType) {
		return _dataValueWrapper
	}
	return _valueWrapper
}

func (dataValueWrapper) Wrap(v reflect.Value) data.Value {
	return v.Interface().(data.Value)
}

func (dataValueWrapper) Unwrap(v data.Value) reflect.Value {
	return reflect.ValueOf(v)
}

func (valueWrapper) Wrap(v reflect.Value) data.Value {
	if w, ok := v.Interface().(*wrappedValue); ok {
		return w
	}
	fmt.Println(v)
	return &wrappedValue{
		elem: v,
	}
}

func (valueWrapper) Unwrap(v data.Value) reflect.Value {
	if w, ok := v.(*wrappedValue); ok {
		return w.elem
	}
	return reflect.ValueOf(v)
}

func (*wrappedValue) Type() data.Name {
	return wrappedValueType
}

func (v *wrappedValue) String() string {
	return data.DumpString(v)
}
