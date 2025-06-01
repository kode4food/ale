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
	bigUintWrapper[T ~uint | ~uint64 | ~uintptr]   struct{}
	smallUintWrapper[T ~uint8 | ~uint16 | ~uint32] struct{}
)

const (
	ErrValueMustBePositiveInteger = "value must be a positive integer"
	ErrValueMustBe64BitInteger    = "value must be a 64-bit integer"
)

func makeWrappedUnsignedInt(t reflect.Type) Wrapper {
	switch k := t.Kind(); k {
	case reflect.Uint:
		return bigUintWrapper[uint]{}
	case reflect.Uintptr:
		return bigUintWrapper[uintptr]{}
	case reflect.Uint64:
		return bigUintWrapper[uint64]{}
	case reflect.Uint32:
		return smallUintWrapper[uint32]{}
	case reflect.Uint16:
		return smallUintWrapper[uint16]{}
	case reflect.Uint8:
		return smallUintWrapper[uint8]{}
	default:
		panic(debug.ProgrammerError("uint kind is incorrect"))
	}
}

func (smallUintWrapper[T]) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Uint()), nil
}

func (smallUintWrapper[T]) Unwrap(v data.Value) (reflect.Value, error) {
	if v, ok := v.(data.Integer); ok {
		return reflect.ValueOf(T(v)), nil
	}
	return zero[T](), errors.New(ErrValueMustBeInteger)
}

func (bigUintWrapper[T]) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	u := v.Uint()
	if u <= math.MaxInt64 {
		return data.Integer(u), nil
	}
	bi := new(big.Int).SetUint64(u)
	return (*data.BigInt)(bi), nil
}

func (bigUintWrapper[T]) Unwrap(v data.Value) (reflect.Value, error) {
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
