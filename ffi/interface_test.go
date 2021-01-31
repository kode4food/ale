package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

type (
	testInterface interface {
		Void(func(testInterface))
		Add(int, int) int
		Double(int, int) (int, int)
		notExported()
	}

	testReceiver bool
)

func testWrap(as *assert.Wrapper) data.Object {
	f := ffi.MustWrap(func() testInterface {
		return testReceiver(true)
	}).(data.Function)
	r := f.Call().(data.Object)
	as.Equal(4, len(r))
	return r
}

func TestNotExported(t *testing.T) {
	as := assert.New(t)
	r := testWrap(as)
	_, ok := r[K("notExported")]
	as.False(ok)
}

func TestVoidInterface(t *testing.T) {
	as := assert.New(t)
	r := testWrap(as)
	b := []bool{false}
	m := r[K("Void")].(data.Function)
	m.Call(ffi.MustWrap(func(_ testInterface) {
		b[0] = true
	}))
	as.True(b[0])
}

func TestInterfaceReceiver(t *testing.T) {
	as := assert.New(t)
	r := testWrap(as)
	b := []bool{false}
	m := r[K("Void")].(data.Function)
	m.Call(ffi.MustWrap(func(r testInterface) {
		r, ok := r.(testReceiver)
		as.True(ok)
		b[0] = bool(r.(testReceiver))
	}))
	as.True(b[0])
}

func TestValueInterface(t *testing.T) {
	as := assert.New(t)
	r := testWrap(as)
	m := r[K("Add")].(data.Function)
	s := m.Call(ffi.MustWrap(I(4)), ffi.MustWrap(I(6)))
	as.Equal(I(10), s)
}

func TestVectorInterface(t *testing.T) {
	as := assert.New(t)
	r := testWrap(as)
	m := r[K("Double")].(data.Function)
	d := m.Call(ffi.MustWrap(I(4)), ffi.MustWrap(I(6))).(data.Vector)
	as.Equal(2, len(d))
	as.Equal(I(8), d[0])
	as.Equal(I(12), d[1])
}

func (r testReceiver) Void(f func(t testInterface)) {
	f(r)
}

func (testReceiver) Add(l, r int) int {
	return l + r
}

func (testReceiver) Double(f, s int) (int, int) {
	return f * 2, s * 2
}

func (testReceiver) notExported() {}
