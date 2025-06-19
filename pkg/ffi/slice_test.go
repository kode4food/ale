package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/ffi"
	"github.com/kode4food/ale/pkg/read"
)

func TestSliceWrap(t *testing.T) {
	as := assert.New(t)
	a1 := []int{1, 2, 3}
	d1 := ffi.MustWrap(a1).(data.Vector)
	as.Equal(3, len(d1))
	as.Equal(I(1), d1[0])
	as.Equal(I(2), d1[1])
	as.Equal(I(3), d1[2])
}

func TestSliceUnwrap(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(a []int) []int {
		return basics.Map(a, func(i int) int {
			return i * 2
		})
	}).(data.Procedure)

	out := f.Call(V(I(1), I(2), I(3))).(data.Vector)
	as.Equal(3, len(out))
	as.Equal(I(2), out[0])
	as.Equal(I(4), out[1])
	as.Equal(I(6), out[2])

	out = f.Call(L(I(3), I(7), I(11), I(25))).(data.Vector)
	as.Equal(4, len(out))
	as.Equal(I(6), out[0])
	as.Equal(I(14), out[1])
	as.Equal(I(22), out[2])
	as.Equal(I(50), out[3])

	ns := assert.GetTestNamespace()
	out = f.Call(read.MustFromString(ns, "5 8 13")).(data.Vector)
	as.Equal(3, len(out))
	as.Equal(I(10), out[0])
	as.Equal(I(16), out[1])
	as.Equal(I(26), out[2])
}

func TestByteSliceUnwrap(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(a []byte) []byte {
		return basics.Map(a, func(i byte) byte {
			return i + 1
		})
	}).(data.Procedure)

	out := f.Call(data.Bytes{1, 2, 3}).(data.Bytes)
	as.Equal(3, len(out))
	as.Equal(byte(2), out[0])
	as.Equal(byte(3), out[1])
	as.Equal(byte(4), out[2])
}
