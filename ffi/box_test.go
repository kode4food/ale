package ffi_test

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
)

func TestUintptrBoxing(t *testing.T) {
	as := assert.New(t)

	four2 := uintptr(42)
	w, err := ffi.WrapType(reflect.TypeOf(four2))
	if as.NoError(err) {
		b1, err := w.Wrap(nil, reflect.ValueOf(four2))
		if as.NoError(err) {
			as.Equal(uint64(42), b1.(*ffi.Boxed[uintptr]).Value.Uint())
		}

		u, err := w.Unwrap(b1)
		if as.NoError(err) {
			as.Equal(reflect.ValueOf(four2), u)
		}

		b2, err := w.Wrap(nil, reflect.ValueOf(uintptr(42)))
		if as.NoError(err) {
			as.True(b1.(*ffi.Boxed[uintptr]).Equal(b2))
		}
	}
}

func TestUnsafePointerBoxing(t *testing.T) {
	as := assert.New(t)

	four2 := 42
	p := unsafe.Pointer(&four2)
	w, err := ffi.WrapType(reflect.TypeOf(p))
	if as.NoError(err) {
		b1, err := w.Wrap(nil, reflect.ValueOf(p))
		if as.NoError(err) {
			as.Equal(
				unsafe.Pointer(&four2),
				b1.(*ffi.Boxed[unsafe.Pointer]).Value.UnsafePointer(),
			)
		}

		u, err := w.Unwrap(b1)
		if as.NoError(err) {
			as.Equal(reflect.ValueOf(p), u)
		}

		b2, err := w.Wrap(nil, reflect.ValueOf(unsafe.Pointer(&four2)))
		if as.NoError(err) {
			as.True(b1.(*ffi.Boxed[unsafe.Pointer]).Equal(b2))
		}
	}
}
