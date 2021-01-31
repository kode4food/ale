package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestBoolWrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(b bool) bool {
		return !b
	}).(data.Function)
	b := f.Call(data.False)
	as.True(b)
}

func TestFloatWrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(f1 float32, f2 float64) (float32, float64) {
		return f1 * 2, f2 * 3
	}).(data.Function)
	r := f.Call(F(9), F(15)).(data.Vector)
	as.Equal(F(18), r[0])
	as.Equal(F(45), r[1])
}
