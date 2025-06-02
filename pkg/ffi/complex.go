package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/pkg/data"
)

type complexWrapper[T ~complex128 | ~complex64] struct{}

const (
	// ErrValueMustBeCons is raised when a complex Unwrap call can't treat its
	// source as a data.Cons
	ErrValueMustBeCons = "value must be a cons cell"

	// ErrConsMustContainFloat is raised when a complex Unwrap call can't treat
	// its source's components as data.Floats
	ErrConsMustContainFloat = "components must be float values"
)

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
