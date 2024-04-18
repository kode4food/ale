package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/ffi"
	"github.com/kode4food/comb/basics"
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
	out := f.Call(data.NewVector(I(1), I(2), I(3))).(data.Vector)
	as.Equal(3, len(out))
	as.Equal(I(2), out[0])
	as.Equal(I(4), out[1])
	as.Equal(I(6), out[2])
}
