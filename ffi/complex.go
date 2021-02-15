package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	complex128Wrapper reflect.Kind
	complex64Wrapper  reflect.Kind
)

// Error messages
const (
	ErrValueMustBeCons      = "value must be a cons cell"
	ErrConsMustContainFloat = "components must be float values"

	errIncorrectComplexKind = "complex kind is incorrect"
)

var (
	complex128zero = reflect.ValueOf(0 + 0i)
	complex64zero  = reflect.ValueOf(complex64(0 + 0i))
)

func makeWrappedComplex(t reflect.Type) (Wrapper, error) {
	k := t.Kind()
	switch k {
	case reflect.Complex128:
		return complex128Wrapper(k), nil
	case reflect.Complex64:
		return complex64Wrapper(k), nil
	default:
		// Programmer error
		panic(errors.New(errIncorrectComplexKind))
	}
}

func (w complex128Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	c := v.Complex()
	r := data.Float(real(c))
	i := data.Float(imag(c))
	return data.NewCons(r, i), nil
}

func (w complex128Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if c, ok := v.(data.Cons); ok {
		r, rok := c.Car().(data.Float)
		i, iok := c.Cdr().(data.Float)
		if rok && iok {
			out := complex(float64(r), float64(i))
			return reflect.ValueOf(out), nil
		}
		return complex128zero, errors.New(ErrConsMustContainFloat)
	}
	return complex128zero, errors.New(ErrValueMustBeCons)
}

func (w complex64Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	c := v.Complex()
	r := data.Float(real(c))
	i := data.Float(imag(c))
	return data.NewCons(r, i), nil
}

func (w complex64Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if c, ok := v.(data.Cons); ok {
		r, rok := c.Car().(data.Float)
		i, iok := c.Cdr().(data.Float)
		if rok && iok {
			out := (complex64)(complex(float64(r), float64(i)))
			return reflect.ValueOf(out), nil
		}
		return complex64zero, errors.New(ErrConsMustContainFloat)
	}
	return complex64zero, errors.New(ErrValueMustBeCons)
}
