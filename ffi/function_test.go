package ffi_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestVoidResult(t *testing.T) {
	var b bool
	as := assert.New(t)
	f := ffi.MustWrap(func(i int) {
		b = i == 37
	}).(data.Function)
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
	}).(data.Function)
	as.NotNil(f)
	r := f.Call(I(5))
	as.Equal(I(10), r)
}

func TestVectorResult(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i int, s string) (int, string) {
		return i * 2, s + "-modified"
	}).(data.Function)
	as.NotNil(f)
	r := f.Call(I(4), S("hello")).(data.Vector)
	as.Equal(I(8), r[0])
	as.Equal(S("hello-modified"), r[1])
}

func TestFuncUnwrap(t *testing.T) {
	as := assert.New(t)
	var set bool
	mark := func() {
		set = true
	}
	f := ffi.MustWrap(func(f func()) func() {
		f()
		as.True(set)
		return f
	}).(data.Function)
	inFunc := ffi.MustWrap(mark)
	res := f.Call(inFunc)
	as.NotNil(res)
	fmt.Println(res)
	as.Contains(":type applicative", res)
}
