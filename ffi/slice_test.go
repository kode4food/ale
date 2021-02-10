package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestSliceWrap(t *testing.T) {
	as := assert.New(t)
	a1 := []int{1, 2, 3}
	d1 := ffi.MustWrap(a1).(data.Vector).Values()
	as.Equal(3, len(d1))
	as.Equal(I(1), d1[0])
	as.Equal(I(2), d1[1])
	as.Equal(I(3), d1[2])
}

func TestSliceUnwrap(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(a []int) []int {
		res := make([]int, len(a))
		for i, v := range a {
			res[i] = v * 2
		}
		return res
	}).(data.Function)
	out := f.Call(data.NewVector(I(1), I(2), I(3))).(data.Vector).Values()
	as.Equal(3, len(out))
	as.Equal(I(2), out[0])
	as.Equal(I(4), out[1])
	as.Equal(I(6), out[2])
}
