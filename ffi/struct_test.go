package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

type stateInfo struct {
	Loop        *stateInfo
	Name        string
	notExported string
	Population  int `ale:"pop"`
}

func testStructStateInfo() *stateInfo {
	return &stateInfo{
		Name:        "California",
		Population:  40,
		notExported: "hello",
	}
}

func TestStructWrap(t *testing.T) {
	as := assert.New(t)
	si := testStructStateInfo()
	m := ffi.MustWrap(si).(*data.Object)
	as.Equal(S("California"), as.MustGet(m, K("Name")))
	as.Equal(I(40), as.MustGet(m, K("pop")))
	_, ok := m.Get(K("notExported"))
	as.False(ok)
}

func TestStructCycle(t *testing.T) {
	as := assert.New(t)
	si := testStructStateInfo()
	si.Loop = si

	res, err := ffi.Wrap(si)
	if as.Nil(res) && as.NotNil(err) {
		as.Equal(ffi.ErrCycleDetected, err.Error())
	}
}

func TestStructUnwrap(t *testing.T) {
	as := assert.New(t)
	si := testStructStateInfo()
	f := ffi.MustWrap(func(i *stateInfo) (string, int) {
		return i.Name, i.Population
	}).(data.Procedure)
	r := f.Call(ffi.MustWrap(si)).(data.Vector)
	as.Equal(S("California"), r[0])
	as.Equal(I(40), r[1])
}

func BenchmarkStructWrapper(b *testing.B) {
	f := ffi.MustWrap(func(i *stateInfo) (string, int) {
		return i.Name, i.Population
	}).(data.Procedure)

	b.ResetTimer()
	for range b.N {
		si := testStructStateInfo()
		r := f.Call(ffi.MustWrap(si)).(data.Vector)
		if S("California") != r[0] || I(40) != r[1] {
			b.Fail()
		}
	}
}
