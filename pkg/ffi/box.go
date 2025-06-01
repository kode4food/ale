package ffi

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/data"
)

type (
	// Boxed is a wrapper for boxed values, which are typically used to
	// represent values that are not directly addressable in Go, such as
	// pointers or uintptrs. It holds a reflect.Value that contains the boxed
	// value
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

func makeBoxedWrapper(t reflect.Type) Wrapper {
	switch k := t.Kind(); k {
	case reflect.Uintptr:
		return boxedWrapper[uintptr]{}
	case reflect.UnsafePointer:
		return boxedWrapper[unsafe.Pointer]{}
	default:
		panic(debug.ProgrammerError("boxed kind is incorrect"))
	}
}

func (b boxedWrapper[T]) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return &Boxed[T]{Value: v}, nil
}

func (b boxedWrapper[T]) Unwrap(v data.Value) (reflect.Value, error) {
	if box, ok := v.(*Boxed[T]); ok {
		return reflect.ValueOf(box.Value), nil
	}
	return zero[reflect.Value](), errors.New(ErrValueMustBeBoxed)
}

func (b *Boxed[T]) Equal(other data.Value) bool {
	if o, ok := other.(*Boxed[T]); ok {
		return b == o || b.Value == o.Value
	}
	return false
}
