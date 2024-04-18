package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/ffi"
)

func TestVoidResult(t *testing.T) {
	var b bool
	as := assert.New(t)
	f := ffi.MustWrap(func(i int) {
		b = i == 37
	}).(data.Procedure)
	as.NotNil(f)
	as.False(b)
	r := f.Call(I(37))
	as.Nil(r)
	as.True(b)
}

func TestSingleResult(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i int) int {
		return i * 2
	}).(data.Procedure)
	as.NotNil(f)
	r := f.Call(I(5))
	as.Equal(I(10), r)
}

func TestVectorResult(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i int, s string) (int, string) {
		return i * 2, s + "-modified"
	}).(data.Procedure)
	as.NotNil(f)
	r := f.Call(I(4), S("hello")).(data.Vector)
	as.Equal(I(8), r[0])
	as.Equal(S("hello-modified"), r[1])
}

func TestVoidFuncUnwrap(t *testing.T) {
	as := assert.New(t)
	var set bool
	mark := func() {
		set = true
	}
	f := ffi.MustWrap(func(f func()) func() {
		f()
		as.True(set)
		return f
	}).(data.Procedure)
	inFunc := ffi.MustWrap(mark)
	res := f.Call(inFunc)
	as.NotNil(res)
	as.Contains(":type procedure", res)
}

func TestValueFuncUnwrap(t *testing.T) {
	as := assert.New(t)
	fourTwo := func() int {
		return 42
	}
	f := ffi.MustWrap(func(f func() int) func() int {
		as.Number(42, f())
		return f
	}).(data.Procedure)
	inFunc := ffi.MustWrap(fourTwo)
	as.NotNil(f.Call(inFunc))
}

func TestVectorFuncUnwrap(t *testing.T) {
	as := assert.New(t)
	hello := func() (int, string) {
		return 42, "hello"
	}
	f := ffi.MustWrap(func(f func() (int, string)) func() (int, string) {
		n, s := f()
		as.Number(42, n)
		as.String("hello", s)
		return f
	}).(data.Procedure)
	inFunc := ffi.MustWrap(hello)
	as.NotNil(f.Call(inFunc))
}
