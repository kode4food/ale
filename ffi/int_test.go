package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestInt64Wrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i1 int32, i2 int64) (int32, int64) {
		return i1 * 2, i2 * 3
	}).(data.Function)
	r := f.Call(I(9), I(15)).(data.Vector).Values()
	as.Equal(I(18), r[0])
	as.Equal(I(45), r[1])
}

func TestInt16Wrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i1 int16, i2 int8) (int16, int8) {
		return i1 * 2, i2 * 3
	}).(data.Function)
	r := f.Call(I(9), I(15)).(data.Vector).Values()
	as.Equal(I(18), r[0])
	as.Equal(I(45), r[1])
}
