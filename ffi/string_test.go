package ffi_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestStringWrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(s string) string {
		return fmt.Sprintf("Hello, %s!", s)
	}).(data.Function)
	s := f.Call(S("Ale"))
	as.String("Hello, Ale!", s)
}
