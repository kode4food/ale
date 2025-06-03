package ffi

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"unsafe"

	"github.com/kode4food/ale/pkg/data"
)

type (
	intWrapper[T ~int | ~int8 | ~int16 | ~int32 | ~int64]       struct{}
	uintWrapper[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64] struct{}
	floatWrapper[T ~float32 | ~float64]                         struct{}
	complexWrapper[T ~complex128 | ~complex64]                  struct{}
)

const (
	// ErrValueMustBeSignedInteger is raised when an int Unwrap call can't
	// properly size its source as a signed integer
	ErrValueMustBeSignedInteger = "value must be a %d-bit signed integer"

	// ErrValueMustBeUnsignedInteger is raised when a uint Unwrap call
	// can't properly size its source as an unsigned integer
	ErrValueMustBeUnsignedInteger = "value must be a %d-bit unsigned integer"

	// ErrValueMustBeCons is raised when a complex Unwrap call can't treat its
	// source as a data.Cons
	ErrValueMustBeCons = "value must be a cons cell"

	// ErrValueMustBeFloat is raised when a float Unwrap call can't treat its
	// source as a data.Integer or data.Float
	ErrValueMustBeFloat = "value must be a float"

	// ErrConsMustContainFloat is raised when a complex Unwrap call can't treat
	// its source's components as data.Floats
	ErrConsMustContainFloat = "components must be float values"
)

func (intWrapper[_]) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Int()), nil
}

func (intWrapper[T]) Unwrap(v data.Value) (reflect.Value, error) {
	bits := int(unsafe.Sizeof(T(0))) * 8
	switch i := v.(type) {
	case data.Integer:
		res := T(i)
		if data.Integer(res) == i {
			return reflect.ValueOf(res), nil
		}
	case *data.BigInt:
		bi := (*big.Int)(i)
		if bi.BitLen() <= bits-1 {
			return reflect.ValueOf(T(bi.Int64())), nil
		}
	}
	return zero[T](), fmt.Errorf(ErrValueMustBeSignedInteger, bits)
}

func (uintWrapper[_]) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	u := v.Uint()
	if u <= math.MaxInt64 {
		return data.Integer(u), nil
	}
	bi := new(big.Int).SetUint64(u)
	return (*data.BigInt)(bi), nil
}

func (uintWrapper[T]) Unwrap(v data.Value) (reflect.Value, error) {
	bits := int(unsafe.Sizeof(T(0))) * 8
	switch i := v.(type) {
	case data.Integer:
		res := T(i)
		if data.Integer(res) == i {
			return reflect.ValueOf(res), nil
		}
	case *data.BigInt:
		bi := (*big.Int)(i)
		if bi.Sign() >= 0 && bi.BitLen() <= bits {
			return reflect.ValueOf(T(bi.Uint64())), nil
		}
	}
	return zero[T](), fmt.Errorf(ErrValueMustBeUnsignedInteger, bits)
}

func (floatWrapper[_]) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Float(v.Float()), nil
}

func (floatWrapper[T]) Unwrap(v data.Value) (reflect.Value, error) {
	if f, ok := makeFloat64(v); ok {
		return reflect.ValueOf(T(f)), nil
	}
	return reflect.Value{}, errors.New(ErrValueMustBeFloat)
}

func makeFloat64(v data.Value) (float64, bool) {
	switch v := v.(type) {
	case data.Integer:
		return float64(v), true
	case data.Float:
		return float64(v), true
	default:
		return 0, false
	}
}

func (complexWrapper[_]) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	c := v.Complex()
	r := data.Float(real(c))
	i := data.Float(imag(c))
	return data.NewCons(r, i), nil
}

func (complexWrapper[T]) Unwrap(v data.Value) (reflect.Value, error) {
	if c, ok := v.(*data.Cons); ok {
		r, rok := c.Car().(data.Float)
		i, iok := c.Cdr().(data.Float)
		if rok && iok {
			out := (T)(complex(r, i))
			return reflect.ValueOf(out), nil
		}
		return zero[T](), errors.New(ErrConsMustContainFloat)
	}
	return zero[T](), errors.New(ErrValueMustBeCons)
}
