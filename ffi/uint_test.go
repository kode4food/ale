package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestUInt64Wrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i1 uint32, i2 uint64) (uint32, uint64) {
		return i1 * 2, i2 * 3
	}).(data.Function)
	r := f.Call(I(9), I(15)).(data.Vector).Values()
	as.Equal(I(18), r[0])
	as.Equal(I(45), r[1])
}

func TestUInt16Wrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i1 uint16, i2 uint8) (uint16, uint8) {
		return i1 * 2, i2 * 3
	}).(data.Function)
	r := f.Call(I(9), I(15)).(data.Vector).Values()
	as.Equal(I(18), r[0])
	as.Equal(I(45), r[1])
}
