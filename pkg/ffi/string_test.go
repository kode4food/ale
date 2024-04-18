package ffi_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/ffi"
)

func TestStringWrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(s string) string {
		return fmt.Sprintf("Hello, %s!", s)
	}).(data.Procedure)
	s := f.Call(S("Ale"))
	as.String("Hello, Ale!", s)
}
