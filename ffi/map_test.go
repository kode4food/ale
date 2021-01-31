package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

type cycleMap map[string]interface{}

var stateMap = map[string]int{
	"California":    40,
	"Massachusetts": 7,
	"Virginia":      8,
}

func TestMapWrap(t *testing.T) {
	as := assert.New(t)
	m := ffi.MustWrap(stateMap).(data.Object)
	as.NotNil(m)
	as.Equal(I(40), m[S("California")])
	as.Equal(I(7), m[S("Massachusetts")])
}

func TestMapCycle(t *testing.T) {
	as := assert.New(t)
	m := cycleMap{
		"k1": 99,
		"k2": 100,
	}
	m["k3"] = m

	res, err := ffi.Wrap(m)
	as.Nil(res)
	as.NotNil(err)
	as.Equal(ffi.ErrCycleDetected, err.Error())
}

func TestMapUnwrap(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(k string, m map[string]int) int {
		return m[k]
	}).(data.Function)
	m := ffi.MustWrap(stateMap).(data.Object)
	as.Equal(I(40), f.Call(S("California"), m))
	as.Equal(I(8), f.Call(S("Virginia"), m))
}
