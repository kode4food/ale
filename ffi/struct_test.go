package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

type stateInfo struct {
	Name       string
	Population int
}

func TestStructWrap(t *testing.T) {
	as := assert.New(t)
	m := ffi.Wrap(&stateInfo{
		Name:       "California",
		Population: 40,
	}).(data.Object)
	as.Equal(S("California"), m[K("Name")])
	as.Equal(I(40), m[K("Population")])
}

func TestStructUnwrap(t *testing.T) {
	as := assert.New(t)
	f := makeCall(ffi.Wrap(func(i *stateInfo) (string, int) {
		return i.Name, i.Population
	}))
	r := f(ffi.Wrap(&stateInfo{
		Name:       "California",
		Population: 40,
	})).(data.Vector)
	as.Equal(S("California"), r[0])
	as.Equal(I(40), r[1])
}
