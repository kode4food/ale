package ffi

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/kode4food/ale"
)

type (
	// Boxed is a wrapper for reflected values that are not addressable in Go,
	// such as unsafe pointers or uintptrs
	Boxed[T boxedTypes] struct {
		Value reflect.Value
	}

	boxedWrapper[T boxedTypes] struct{}

	boxedTypes interface {
		~uintptr | unsafe.Pointer
	}
)

const (
	ErrValueMustBeBoxed = "value must be a boxed value"
)

func (b boxedWrapper[T]) Wrap(_ *Context, v reflect.Value) (ale.Value, error) {
	return &Boxed[T]{Value: v}, nil
}

func (b boxedWrapper[T]) Unwrap(v ale.Value) (reflect.Value, error) {
	if box, ok := v.(*Boxed[T]); ok {
		return box.Value, nil
	}
	return _zero, errors.New(ErrValueMustBeBoxed)
}

func (b *Boxed[T]) Equal(other ale.Value) bool {
	if o, ok := other.(*Boxed[T]); ok {
		return b == o || b.Value == o.Value || b.Value.Equal(o.Value)
	}
	return false
}
