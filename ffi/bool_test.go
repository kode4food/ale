package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
)

func TestBoolWrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(b bool) bool {
		return !b
	}).(data.Function)

	b := f.Call(data.False)
	as.True(b)

	b = f.Call(data.True)
	as.False(b)
}
