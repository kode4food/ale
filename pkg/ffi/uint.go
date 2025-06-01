package ffi

import (
	"errors"
	"math"
	"math/big"
	"reflect"

	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/data"
)

type (
	uintWrapper    reflect.Kind
	uint64Wrapper  reflect.Kind
	uintptrWrapper reflect.Kind
	uint32Wrapper  reflect.Kind
	uint16Wrapper  reflect.Kind
	uint8Wrapper   reflect.Kind

	unwrappableInts interface {
		~uint | ~uint64 | ~uintptr
	}
)

const (
	ErrValueMustBePositiveInteger = "value must be a positive integer"
	ErrValueMustBe64BitInteger    = "value must be a 64-bit integer"
)

func makeWrappedUnsignedInt(t reflect.Type) Wrapper {
	switch k := t.Kind(); k {
	case reflect.Uint:
		return uintWrapper(k)
	case reflect.Uintptr:
		return uintptrWrapper(k)
	case reflect.Uint64:
		return uint64Wrapper(k)
	case reflect.Uint32:
		return uint32Wrapper(k)
	case reflect.Uint16:
		return uint16Wrapper(k)
	case reflect.Uint8:
		return uint8Wrapper(k)
	default:
		panic(debug.ProgrammerError("uint kind is incorrect"))
	}
}

func (uintWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return wrapUint64(v.Uint()), nil
}

func (uintWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapUint64[uint](v)
}

func (uintptrWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return wrapUint64(v.Uint()), nil
}

func (uintptrWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapUint64[uintptr](v)
}

func (uint64Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return wrapUint64(v.Uint()), nil
}

func (uint64Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapUint64[uint64](v)
}

func (uint32Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Uint()), nil
}

func (uint32Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapInt[uint32](v)
}

func (uint16Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Uint()), nil
}

func (uint16Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapInt[uint16](v)
}

func (uint8Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Uint()), nil
}

func (uint8Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapInt[uint8](v)
}

func wrapUint64(u uint64) data.Value {
	if u <= math.MaxInt64 {
		return data.Integer(u)
	}
	bi := new(big.Int).SetUint64(u)
	return (*data.BigInt)(bi)
}

func unwrapUint64[T unwrappableInts](v data.Value) (reflect.Value, error) {
	switch i := v.(type) {
	case data.Integer:
		if i < 0 {
			return zero[T](), errors.New(ErrValueMustBePositiveInteger)
		}
		return reflect.ValueOf(T(uint64(i))), nil
	case *data.BigInt:
		bi := (*big.Int)(i)
		if bi.Sign() < 0 {
			return zero[T](), errors.New(ErrValueMustBePositiveInteger)
		}
		if bi.BitLen() > 64 {
			return zero[T](), errors.New(ErrValueMustBe64BitInteger)
		}
		return reflect.ValueOf(T(bi.Uint64())), nil
	}
	return zero[T](), errors.New(ErrValueMustBeInteger)
}
