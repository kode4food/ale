package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/ffi"
)

func TestIntWrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i1 int, i2 int) int {
		return i1 + i2
	}).(data.Procedure)
	r := f.Call(I(9), I(15))
	as.Equal(I(24), r)
}

func TestInt64Wrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i1 int32, i2 int64) (int32, int64) {
		return i1 * 2, i2 * 3
	}).(data.Procedure)
	r := f.Call(I(9), I(15)).(data.Vector)
	as.Equal(I(18), r[0])
	as.Equal(I(45), r[1])
}

func TestInt16Wrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i1 int16, i2 int8) (int16, int8) {
		return i1 * 2, i2 * 3
	}).(data.Procedure)
	r := f.Call(I(9), I(50)).(data.Vector)
	as.Equal(I(18), r[0])
	as.Equal(I(-106), r[1])
}

func TestUIntWrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i1 uint, i2 uint) uint {
		return i1 + i2
	}).(data.Procedure)
	r := f.Call(I(9), I(15))
	as.Equal(I(24), r)
}

func TestUInt64Wrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i1 uint32, i2 uint64) (uint32, uint64) {
		return i1 * 2, i2 * 3
	}).(data.Procedure)
	r := f.Call(I(9), I(15)).(data.Vector)
	as.Equal(I(18), r[0])
	as.Equal(I(45), r[1])
}

func TestUInt16Wrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i1 uint16, i2 uint8) (uint16, uint8) {
		return i1 * 2, i2 * 3
	}).(data.Procedure)
	r := f.Call(I(9), I(15)).(data.Vector)
	as.Equal(I(18), r[0])
	as.Equal(I(45), r[1])
}

func TestFloatWrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(f1 float32, f2 float64) (float32, float64) {
		return f1 * 2, f2 * 3
	}).(data.Procedure)
	r := f.Call(F(9), I(15)).(data.Vector)
	as.Equal(F(18), r[0])
	as.Equal(F(45), r[1])
}

func TestFloatEval(t *testing.T) {
	as := NewWrapped(t)

	as.EvalTo(
		`(d 2.5 2.4)`,
		Env{
			"d": func(f32 float32, f64 float64) (float32, float64) {
				return f32 * 2, f64 * 2
			},
		},
		V(F(5.0), F(4.8)),
	)
}

func TestComplexWrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(
		func(i1 complex64, i2 complex128) (complex64, complex128) {
			return i1 * 2, i2 * 3
		},
	).(data.Procedure)
	c1 := C(F(9), F(15))
	c2 := C(F(32), F(2))
	r := f.Call(c1, c2).(data.Vector)
	as.String("(18.0 . 30.0)", r[0])
	as.String("(96.0 . 6.0)", r[1])
}
