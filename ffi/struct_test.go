package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

type stateInfo struct {
	Name        string
	Population  int `ale:"pop"`
	Loop        *stateInfo
	notExported string
}

func TestStructWrap(t *testing.T) {
	as := assert.New(t)
	si := &stateInfo{
		Name:        "California",
		Population:  40,
		notExported: "hello",
	}
	m := ffi.MustWrap(si).(data.Object)
	as.Equal(S("California"), m[K("Name")])
	as.Equal(I(40), m[K("pop")])
	_, ok := m[K("notExported")]
	as.False(ok)
}

func TestStructCycle(t *testing.T) {
	as := assert.New(t)
	si := &stateInfo{
		Name:        "California",
		Population:  40,
		notExported: "hello",
	}
	si.Loop = si

	res, err := ffi.Wrap(si)
	as.Nil(res)
	as.NotNil(err)
	as.Equal(ffi.ErrCycleDetected, err.Error())
}

func TestStructUnwrap(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(i *stateInfo) (string, int) {
		return i.Name, i.Population
	}).(data.Function)
	r := f.Call(ffi.MustWrap(&stateInfo{
		Name:       "California",
		Population: 40,
	})).(data.Vector)
	as.Equal(S("California"), r[0])
	as.Equal(I(40), r[1])
}
