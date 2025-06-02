package ffi

import (
	"errors"
	"math"
	"math/big"
	"reflect"

	"github.com/kode4food/ale/pkg/data"
)

type (
	intWrapper[T wrappableInts] struct{}

	wrappableInts interface {
		~int | ~uint8 | ~int8 | ~uint16 | ~int16 | ~uint32 | ~int32 | ~int64
	}

	uint64Wrapper[T ~uint | ~uint64]           struct{}
	uintWrapper[T ~uint8 | ~uint16 | ~uint32]  struct{}
	floatWrapper[T ~float32 | ~float64]        struct{}
	complexWrapper[T ~complex128 | ~complex64] struct{}
)

const (
	// ErrValueMustBeInteger is raised when an integer Unwrap call can't treat
	// its source as a data.Integer
	ErrValueMustBeInteger = "value must be an integer"

	ErrValueMustBePositiveInteger = "value must be a positive integer"
	ErrValueMustBe64BitInteger    = "value must be a 64-bit integer"

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
	if v, ok := v.(data.Integer); ok {
		return reflect.ValueOf(T(v)), nil
	}
	return zero[T](), errors.New(ErrValueMustBeInteger)
}

func (uintWrapper[_]) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Uint()), nil
}

func (uintWrapper[T]) Unwrap(v data.Value) (reflect.Value, error) {
	if v, ok := v.(data.Integer); ok {
		return reflect.ValueOf(T(v)), nil
	}
	return zero[T](), errors.New(ErrValueMustBeInteger)
}

func (uint64Wrapper[_]) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	u := v.Uint()
	if u <= math.MaxInt64 {
		return data.Integer(u), nil
	}
	bi := new(big.Int).SetUint64(u)
	return (*data.BigInt)(bi), nil
}

func (uint64Wrapper[T]) Unwrap(v data.Value) (reflect.Value, error) {
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
