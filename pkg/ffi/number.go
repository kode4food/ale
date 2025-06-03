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
	intWrapper[T intType]   struct{}
	uintWrapper[T uintType] struct{}

	floatWrapper[T ~float32 | ~float64]        struct{}
	complexWrapper[T ~complex128 | ~complex64] struct{}

	intType interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64
	}

	uintType interface {
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
	}
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
	switch i := v.(type) {
	case data.Integer:
		if res, ok := int64ToInt[T](int64(i)); ok {
			return reflect.ValueOf(res), nil
		}
	case data.Float:
		if res, ok := floatToInt[T](float64(i)); ok {
			return reflect.ValueOf(res), nil
		}
	case *data.BigInt:
		if bi := (*big.Int)(i); bi.IsInt64() {
			if res, ok := int64ToInt[T](bi.Int64()); ok {
				return reflect.ValueOf(res), nil
			}
		}
	case *data.Ratio:
		if r := (*big.Rat)(i); r.IsInt() {
			if res, ok := int64ToInt[T](r.Num().Int64()); ok {
				return reflect.ValueOf(res), nil
			}
		}
	}
	bits := int(unsafe.Sizeof(T(0))) * 8
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
	switch i := v.(type) {
	case data.Integer:
		if i >= 0 {
			if res, ok := uint64ToUint[T](uint64(i)); ok {
				return reflect.ValueOf(res), nil
			}
		}
	case data.Float:
		if res, ok := floatToUint[T](float64(i)); ok {
			return reflect.ValueOf(res), nil
		}
	case *data.BigInt:
		if bi := (*big.Int)(i); bi.IsUint64() {
			if res, ok := uint64ToUint[T](bi.Uint64()); ok {
				return reflect.ValueOf(res), nil
			}
		}
	case *data.Ratio:
		if r := (*big.Rat)(i); r.IsInt() {
			if bi := r.Num(); bi.IsUint64() {
				if res, ok := uint64ToUint[T](r.Num().Uint64()); ok {
					return reflect.ValueOf(res), nil
				}
			}
		}
	}
	bits := int(unsafe.Sizeof(T(0))) * 8
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

func (complexWrapper[_]) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	c := v.Complex()
	r := data.Float(real(c))
	i := data.Float(imag(c))
	return data.NewCons(r, i), nil
}

func (complexWrapper[T]) Unwrap(v data.Value) (reflect.Value, error) {
	if c, ok := v.(*data.Cons); ok {
		r, rok := makeFloat64(c.Car())
		i, iok := makeFloat64(c.Cdr())
		if rok && iok {
			out := (T)(complex(r, i))
			return reflect.ValueOf(out), nil
		}
		return zero[T](), errors.New(ErrConsMustContainFloat)
	}
	return zero[T](), errors.New(ErrValueMustBeCons)
}

func makeFloat64(v data.Value) (float64, bool) {
	switch v := v.(type) {
	case data.Integer:
		return float64(v), true
	case data.Float:
		return float64(v), true
	case *data.BigInt:
		f, _ := (*big.Int)(v).Float64()
		return f, true
	case *data.Ratio:
		f, _ := (*big.Rat)(v).Float64()
		return f, true
	default:
		return 0, false
	}
}

func int64ToInt[T intType](i int64) (T, bool) {
	if res := T(i); int64(res) == i {
		return res, true
	}
	return T(0), false
}

func uint64ToUint[T uintType](i uint64) (T, bool) {
	if res := T(i); uint64(res) == i {
		return res, true
	}
	return T(0), false
}

func floatToInt[T intType](f float64) (T, bool) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return T(0), false
	}
	w, r := math.Modf(f)
	if r != 0 {
		return T(0), false
	}
	if res := T(f); int64(res) == int64(w) {
		return res, true
	}
	return T(0), false
}

func floatToUint[T uintType](f float64) (T, bool) {
	if f < 0 || math.IsNaN(f) || math.IsInf(f, 0) {
		return T(0), false
	}
	w, r := math.Modf(f)
	if r != 0 || w < 0 {
		return T(0), false
	}
	if res := T(f); uint64(res) == uint64(w) {
		return res, true
	}
	return T(0), false
}
