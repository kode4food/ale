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

var (
	complex128zero = reflect.ValueOf(0 + 0i)
	complex64zero  = reflect.ValueOf(complex64(0 + 0i))
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
	c := v.Complex()
	r := data.Float(real(c))
	i := data.Float(imag(c))
	return data.NewCons(r, i), nil
}

func (complex128Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if c, ok := v.(*data.Cons); ok {
		r, rok := makeFloat64(c.Car())
		i, iok := makeFloat64(c.Cdr())
		if rok && iok {
			out := complex(r, i)
			return reflect.ValueOf(out), nil
		}
		return complex128zero, errors.New(ErrConsMustContainFloat)
	}
	return complex128zero, errors.New(ErrValueMustBeCons)
}

func (complex64Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	c := v.Complex()
	r := data.Float(real(c))
	i := data.Float(imag(c))
	return data.NewCons(r, i), nil
}

func (complex64Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if c, ok := v.(*data.Cons); ok {
		r, rok := c.Car().(data.Float)
		i, iok := c.Cdr().(data.Float)
		if rok && iok {
			out := (complex64)(complex(r, i))
			return reflect.ValueOf(out), nil
		}
		return complex64zero, errors.New(ErrConsMustContainFloat)
	}
	return complex64zero, errors.New(ErrValueMustBeCons)
}
