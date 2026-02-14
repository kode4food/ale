package ffi_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
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

func TestIntWrapperErrorCases(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i int8) int8 {
		return i
	}).(data.Procedure)

	as.Equal(I(99), f.Call(F(99)))

	signedErr := fmt.Errorf(ffi.ErrValueMustBeSigned, 8)
	as.Panics(func() { _ = f.Call(F(99.5)) }, signedErr)
	as.Panics(func() { _ = f.Call(F(math.NaN())) }, signedErr)
	as.Panics(func() { _ = f.Call(F(math.Inf(1))) }, signedErr)
	as.Panics(func() { _ = f.Call(F(128)) }, signedErr)
	as.Panics(func() { _ = f.Call(mustRatio(t, "3/2")) }, signedErr)
	as.Panics(func() { _ = f.Call(mustInteger(t, "100000000000000000000")) },
		signedErr,
	)
}

func TestUintWrapperErrorCases(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i uint8) uint8 {
		return i
	}).(data.Procedure)

	as.Equal(I(12), f.Call(F(12)))

	unsignedErr := fmt.Errorf(ffi.ErrValueMustBeUnsigned, 8)
	as.Panics(func() { _ = f.Call(I(-1)) }, unsignedErr)
	as.Panics(func() { _ = f.Call(F(-1)) }, unsignedErr)
	as.Panics(func() { _ = f.Call(F(12.5)) }, unsignedErr)
	as.Panics(func() { _ = f.Call(F(math.NaN())) }, unsignedErr)
	as.Panics(func() { _ = f.Call(F(math.Inf(1))) }, unsignedErr)
	as.Panics(func() { _ = f.Call(F(256)) }, unsignedErr)
	as.Panics(func() { _ = f.Call(mustRatio(t, "3/2")) }, unsignedErr)
	as.Panics(func() { _ = f.Call(mustInteger(t, "100000000000000000000")) },
		unsignedErr,
	)
}

func TestFloatWrapperAdditionalCases(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(v float64) float64 {
		return v
	}).(data.Procedure)

	as.Equal(F(1.5), f.Call(mustRatio(t, "3/2")))
	out := f.Call(mustInteger(t, "100000000000000000000")).(data.Float)
	as.True(float64(out) > 0)
	as.Panics(func() { _ = f.Call(S("bad-float")) },
		errors.New(ffi.ErrValueMustBeFloat),
	)
}

func TestComplexWrapperErrors(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(v complex64) complex64 {
		return v
	}).(data.Procedure)

	as.Panics(func() { _ = f.Call(S("not-cons")) },
		errors.New(ffi.ErrValueMustBeCons),
	)
	as.Panics(func() { _ = f.Call(C(S("bad"), F(1))) },
		errors.New(ffi.ErrConsMustContainFloat),
	)
}

func mustInteger(t *testing.T, s string) data.Number {
	t.Helper()
	res, err := data.ParseInteger(s)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func mustRatio(t *testing.T, s string) data.Number {
	t.Helper()
	res, err := data.ParseRatio(s)
	if err != nil {
		t.Fatal(err)
	}
	return res
}
