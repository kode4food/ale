package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func makeCall(v data.Value) data.Call {
	return v.(data.Caller).Call()
}

func TestVoidResult(t *testing.T) {
	var b bool
	as := assert.New(t)
	f := makeCall(ffi.Wrap(func(i int) {
		b = i == 37
	}))
	as.NotNil(f)
	as.False(b)
	r := f(I(37))
	as.Nil(r)
	as.True(b)
}

func TestSingleResult(t *testing.T) {
	as := assert.New(t)
	f := makeCall(ffi.Wrap(func(i int) int {
		return i * 2
	}))
	as.NotNil(f)
	r := f(I(5))
	as.Equal(I(10), r)
}

func TestVectorResult(t *testing.T) {
	as := assert.New(t)
	f := makeCall(ffi.Wrap(func(i int, s string) (int, string) {
		return i * 2, s + "-modified"
	}))
	as.NotNil(f)
	r := f(I(4), S("hello")).(data.Vector)
	as.Equal(I(8), r[0])
	as.Equal(S("hello-modified"), r[1])
}

func TestFuncUnwrap(t *testing.T) {
	as := assert.New(t)
	var set bool
	mark := func() {
		set = true
	}
	f := makeCall(ffi.Wrap(func(f func()) func() {
		f()
		as.True(set)
		return f
	}))
	inFunc := ffi.Wrap(mark)
	res := f(inFunc)
	as.NotNil(res)
	as.Contains(":type wrapped-func", res)
}
