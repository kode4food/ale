package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestComplexWrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i1 complex64, i2 complex128) (complex64, complex128) {
		return i1 * 2, i2 * 3
	}).(data.Function)
	c1 := C(F(9), F(15))
	c2 := C(F(32), F(2))
	r := f.Call(c1, c2).(data.Vector).Values()
	as.String("(18 . 30)", r[0])
	as.String("(96 . 6)", r[1])
}
