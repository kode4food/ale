package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/data"
)

type (
	complex128Wrapper reflect.Kind
	complex64Wrapper  reflect.Kind
)

const (
	// ErrValueMustBeCons is raised when a complex Unwrap call can't treat its
	// source as a data.Cons
	ErrValueMustBeCons = "value must be a cons cell"

	// ErrConsMustContainFloat is raised when a complex Unwrap call can't treat
	// its source's components as data.Floats
	ErrConsMustContainFloat = "components must be float values"
)

func makeWrappedComplex(t reflect.Type) Wrapper {
	switch k := t.Kind(); k {
	case reflect.Complex128:
		return complex128Wrapper(k)
	case reflect.Complex64:
		return complex64Wrapper(k)
	default:
		panic(debug.ProgrammerError("complex kind is incorrect"))
	}
}

func (complex128Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return wrapComplex(v.Complex()), nil
}

func (complex128Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapComplex[complex128](v)
}

func (complex64Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return wrapComplex(v.Complex()), nil
}

func (complex64Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	return unwrapComplex[complex64](v)
}

func wrapComplex(c complex128) data.Value {
	r := data.Float(real(c))
	i := data.Float(imag(c))
	return data.NewCons(r, i)
}

func unwrapComplex[T ~complex64 | ~complex128](v data.Value) (reflect.Value, error) {
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
