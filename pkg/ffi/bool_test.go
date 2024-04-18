package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/ffi"
)

func TestBoolWrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(b bool) bool {
		return !b
	}).(data.Procedure)

	b := f.Call(data.False)
	as.True(b)

	b = f.Call(data.True)
	as.False(b)
}
